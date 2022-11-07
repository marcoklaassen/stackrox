import { CollectionRequest, CollectionResponse, SelectorRule } from 'services/CollectionsService';
import { ensureExhaustive } from 'utils/type.utils';
import {
    Collection,
    SelectorField,
    SelectorEntityType,
    isSupportedSelectorField,
    isByNameField,
    isByLabelField,
    selectorEntityTypes,
    isSelectorField,
} from './types';

const fieldToEntityMap: Record<SelectorField, SelectorEntityType> = {
    Deployment: 'Deployment',
    'Deployment Label': 'Deployment',
    'Deployment Annotation': 'Deployment',
    Namespace: 'Namespace',
    'Namespace Label': 'Namespace',
    'Namespace Annotation': 'Namespace',
    Cluster: 'Cluster',
    'Cluster Label': 'Cluster',
    'Cluster Annotation': 'Cluster',
};

const LABEL_SEPARATOR = '=';

/**
 * This function takes a raw `CollectionResponse` from the server and parses it into a representation
 * of a `Collection` that can be supported by the current UI controls. If any incompatibilities are detected
 * it will return a list of validation errors to the caller.
 */
export function parseCollection(data: CollectionResponse): Collection | AggregateError {
    const collection: Collection = {
        name: data.name,
        description: data.description,
        inUse: data.inUse,
        embeddedCollectionIds: data.embeddedCollections.map(({ id }) => id),
        resourceSelector: {
            Deployment: { type: 'All' },
            Namespace: { type: 'All' },
            Cluster: { type: 'All' },
        },
    };

    const errors: string[] = [];

    if (data.resourceSelectors.length > 1) {
        errors.push(
            `Multiple 'ResourceSelectors' were found for this collection. Only a single resource selector is supported in the UI. Further validation errors will only apply to the first resource selector in the response.`
        );
    }

    data.resourceSelectors[0]?.rules.forEach((rule) => {
        const field = rule.fieldName;
        if (!isSelectorField(field)) {
            errors.push(`An invalid field name was detected. Found field name [${field}].`);
            return;
        }
        const entity = fieldToEntityMap[field];
        const selectorScope = collection.resourceSelector[entity];
        const existingEntityField = selectorScope.type === 'All' ? undefined : selectorScope.field;
        const hasMultipleFieldsForEntity = existingEntityField && existingEntityField !== field;
        const isUnsupportedField = !isSupportedSelectorField(field);
        const isUnsupportedRuleOperator = rule.operator !== 'OR';

        if (hasMultipleFieldsForEntity) {
            errors.push(
                `Each entity type can only contain rules for a single field. A new rule was found for [${entity} -> ${field}], when rules have already been applied for [${entity} -> ${existingEntityField}].`
            );
        }
        if (isUnsupportedField) {
            errors.push(
                `Collection rules for 'Annotation' field names are not supported at this time. Found field name [${field}].`
            );
        }
        if (isUnsupportedRuleOperator) {
            errors.push(
                `Only the disjunction operation ('OR') is currently supported in the front end collection editor. Received an operator of [${rule.operator}].`
            );
        }

        if (hasMultipleFieldsForEntity || isUnsupportedField || isUnsupportedRuleOperator) {
            return;
        }

        if (selectorScope.type === 'All') {
            if (isByLabelField(field)) {
                collection.resourceSelector[entity] = {
                    type: 'ByLabel',
                    field,
                    rules: [],
                };
            } else if (isByNameField(field)) {
                collection.resourceSelector[entity] = {
                    type: 'ByName',
                    field,
                    rule: { operator: 'OR', values: [] },
                };
            }
        }

        const selector = collection.resourceSelector[entity];

        switch (selector.type) {
            case 'ByLabel': {
                const firstValue = rule.values[0]?.value;

                if (firstValue && firstValue.includes(LABEL_SEPARATOR)) {
                    const key = firstValue.split(LABEL_SEPARATOR)[0] ?? '';
                    selector.rules.push({
                        operator: 'OR',
                        key,
                        // TODO Verify with BE whether or not this is a valid method to get the label values. Is
                        //      it possible that multiple `=` symbols will appear in the data here?
                        values: rule.values.map(({ value }) => value.split('=')[1] ?? ''),
                    });
                }
                break;
            }
            case 'ByName':
                selector.rule.values = rule.values.map(({ value }) => value);
                break;
            case 'All':
                // Do nothing
                break;
            default:
                ensureExhaustive(selector);
        }
    });

    if (errors.length > 0) {
        return new AggregateError(errors);
    }

    return collection;
}

/**
 * This function takes the UI supported representation of a `Collection` and generates
 * a server response in the generic `CollectionRequest` format supported by the back end.
 */
export function generateRequest(collection: Collection): CollectionRequest {
    const rules: SelectorRule[] = [];

    selectorEntityTypes.forEach((entityType) => {
        const selector = collection.resourceSelector[entityType];

        switch (selector.type) {
            case 'ByName':
                rules.push({
                    fieldName: selector.field,
                    operator: selector.rule.operator,
                    values: selector.rule.values.map((value) => ({ value })),
                });
                break;
            case 'ByLabel':
                selector.rules.forEach(({ operator, key, values }) => {
                    rules.push({
                        fieldName: selector.field,
                        operator,
                        values: values.map((value) => ({
                            value: `${key}${LABEL_SEPARATOR}${value}`,
                        })),
                    });
                });
                break;
            case 'All':
                // Append nothing
                break;
            default:
                ensureExhaustive(selector);
        }
    });

    return {
        name: collection.name,
        description: collection.description,
        embeddedCollectionIds: collection.embeddedCollectionIds,
        resourceSelectors: [{ rules }],
    };
}
