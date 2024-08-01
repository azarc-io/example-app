# load extensions
load('ext://dotenv', 'dotenv')
load('ext://helm_resource', 'helm_resource')
load('ext://dotenv', 'dotenv')
load('ext://restart_process', 'docker_build_with_restart')
load('ext://helm_resource', 'helm_resource')

# .env support
dotenv(fn='.env')

# tilt_config.json support
config.define_bool("hmr")
config.define_string("arch")
config.define_bool("debug")
config.define_bool("gateway")
config.define_string("kube_context")
config.define_string("debug_port")
cfg = config.parse()

# arch type to build for, default is amd64, you can update this in the tilt_config.json file
arch = cfg.get('arch', 'amd64')

# if true gateway will be deployed
gateway = cfg.get('gateway', False)

# if true will enable debugging and forward the debug port
debug = cfg.get('debug', False)
debugPort = cfg.get('debug_port', '40000')

# enabling hmr will cause the gateway to stop proxying the front end directly and instead the ingress
# will point to the front end on the existing host, the user will notice no difference but hot reloading
# will be enabled and the build will run in kubernetes instead of your local machine
hmr = cfg.get('hmr', False)

# what kube context to permit, prevents you from switching your kube context to staging for eg. and mistakenly deploying there
allow_k8s_contexts(cfg.get('kube_context', cfg.get('kube_context', 'k3d-dev-1')))

# default registry configuration
default_registry(
    'k3d-local-registry:5000',
    host_from_cluster='localhost:5000'
)

local_resource(
  'backend-compile',
  'CGO_ENABLED=0 GOOS=linux GOARCH=%s go build -gcflags "all=-N -l" -o ./bin/example-app ./cmd/app/main.go' % arch,
  deps=['./cmd/app/main.go', './cmd/app', './internal/', './pkg/'],
  ignore=['./internal/gql/node_modules'],
  labels=["compile"],
  resource_deps=[]
)

# the default entry path
entrypoint = '/example-app'
appDockerFile = './deployment/docker/app/Dockerfile.tilt'
webDockerFile = './deployment/docker/web/Dockerfile.tilt'

# entry path to use if debug is enabled
if debug:
    entrypoint = '/dlv --listen=:%s --api-version=2 --headless=true --only-same-user=false --accept-multiclient exec --continue /example-app' % debugPort

# dockerfile to use if hmr is enabled
if hmr:
    webDockerFile = './deployment/docker/web/Dockerfile.hmr.tilt'

# watches directories for changes and triggers an update of the docker image
docker_build_with_restart(
  'example-app-be',
  context='.',
  entrypoint=entrypoint,
  dockerfile=appDockerFile,
  platform='linux/%s' % arch,
  only=[
    './bin',
  ],
  live_update=[
    sync('./bin/', '/'),
  ]
)

# hmr mode will sync the source files to the sidecar and the build will happen there, in this mode
# the gateway will not proxy but instead the ingress will point at the side car for / where as anything
# else such as /graphql or /api would point at the gateway
if hmr:
    docker_build(
      'example-app-fe',
      context='.',
      entrypoint='yarn dev:tilt',
      dockerfile=webDockerFile,
      platform='linux/%s' % arch,
      only=[
        './cmd/web',
      ],
      ignore=[
        './cmd/web/node_modules',
        './cmd/web/dist'
      ],
      live_update=[
        fall_back_on(['./cmd/web/package.json', './cmd/web/yarn.lock']),
        sync('./cmd/web', '/src'),
      ]
    )
# non hmr mode will simply sync your local dist folder, you can use vite build --watch
# in this mode the gateway will serve the files and the side car becomes an ephemeral initContainer
else:
    local_resource(
      'frontend-build',
      'yarn build --env-mode tilt',
      dir="./cmd/web",
      deps=['./cmd/web'],
      ignore=['./cmd/web/node_modules', './cmd/web/dist'],
      labels=["compile"],
      resource_deps=[]
    )
    docker_build_with_restart(
      'example-app-fe',
      context='.',
      entrypoint='echo "reloading" && cp -R /static/. /web/ && sleep 9999999999d',
      dockerfile=webDockerFile,
      platform='linux/%s' % arch,
      only=[
        './cmd/web/dist',
      ],
      live_update=[
        sync('./cmd/web/dist', '/static'),
      ]
    )

# watches the chart directory and triggers an update if yaml files change
flags = [
    '--values=./deployment/charts/app/values.yaml',
    '--set=dev=true',
    '--set=bind.debug=%s' % debugPort,
    '--set=hmr=%s' % hmr,
]

if hmr:
    flags.append("--values=./deployment/charts/app/values-hmr.yaml")

helm_resource(
  'example-app-chart',
  './deployment/charts/app',
  namespace=os.getenv('NAMESPACE'),
  deps=["./deployment/charts/app"],
  flags=flags,
  image_deps=['example-app-be', 'example-app-fe'],
  image_keys=[('image.repository', 'image.tag'), ('image.web_repository', 'image.web_tag')],
)

if debug:
    k8s_resource('example-app-chart',
        port_forwards=[
            debugPort,  # debugger
        ],
        labels=['deployment'],
        links=[
            'localhost:%s' % debugPort,
        ]
    )
