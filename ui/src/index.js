import React from 'react';
import ReactDOM from 'react-dom';
// eslint-disable-next-line
import 'index.css';  // this file is generated by tailwind (see package.json scripts)
import AppPage from 'Containers/AppPage';
import registerServiceWorker from 'registerServiceWorker';
import apiInterceptors from 'Providers/apiInterceptors';

// eslint-disable-next-line no-unused-vars
import stringUtil from 'utils/string';

ReactDOM.render(<AppPage />, document.getElementById('root'));

apiInterceptors.addRequestTokenInterceptor();
apiInterceptors.addResponseInterceptor();
registerServiceWorker();
