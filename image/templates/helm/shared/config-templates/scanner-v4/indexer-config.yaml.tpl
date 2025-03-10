{{- /*
  This is the configuration file template for the Scanner v4 Indexer.
  Except for in extremely rare circumstances, you DO NOT need to modify this file.
  All config options that are possibly dynamic are templated out and can be modified
  via `--set`/values-files specified via `-f`.
     */ -}}

# Configuration file for Scanner v4 Indexer.
indexer:
  enable: true
  database:
    conn_string: >
      host=scanner-v4-db.{{ .Release.Namespace }}.svc
      port=5432
      sslrootcert=/run/secrets/stackrox.io/certs/ca.pem
      user=postgres
      sslmode=verify-full
      {{ if not (kindIs "invalid" ._rox.scannerV4.db.source.statementTimeoutMs) -}} statement_timeout={{._rox.scannerV4.db.source.statementTimeoutMs}} {{- end }}
      {{ if not (kindIs "invalid" ._rox.scannerV4.db.source.minConns) -}} pool_min_conns={{._rox.scannerV4.db.source.minConns}} {{- end }}
      {{ if not (kindIs "invalid" ._rox.scannerV4.db.source.maxConns) -}} pool_max_conns={{._rox.scannerV4.db.source.maxConns}} {{- end }}
      client_encoding=UTF8
    password_file: /run/secrets/stackrox.io/secrets/password
  get_layer_timeout: 1m
matcher:
  enable: false
log_level: info
grpc_listen_addr: 0.0.0.0:8443
http_listen_addr: 0.0.0.0:9443
