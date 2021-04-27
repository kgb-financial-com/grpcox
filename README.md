# gRPCox
[![Go Report Card](https://goreportcard.com/badge/github.com/gusaul/grpcox)](https://goreportcard.com/report/github.com/gusaul/grpcox)

turn [gRPCurl](https://github.com/fullstorydev/grpcurl) into web based UI, extremely easy to use

This fork [kgb-financial-com/grpcox](https://github.com/kgb-financial-com/grpcox) extends
the original with the option to pre-define connection targets, and to select a new or
existing target via html parameters. It also comes with a new frontend, written in
[vue.js](https://vuejs.org/).

## Extensions for financial.com

### Integration of vue.js

With node.js installed, vue-cli installed (using `npm install -g @vue/cli`),
and from the project directory, 
the following steps have been performed in order to integrate vue.js:

* Issue the command `vue create grpcox` (choosing vue 3 when asked)
* Remove `README.md` and `.gitignore` from the newly created directory `grpcox`. Add
  entries as required to the project directory `.gitignore`.
* Move all files and directories from `grpcox` one level up, and delete the 
  then-empty directory.
  
The directories `src` and `public` now contain some Vue.js *Hello World* application.
We are rewriting it according to our needs. We also change `routes.go` such that instead
of `/index` directory, our go software will use `/dist` as the source for our html pages.
Now, we can develop the frontend by starting it the usual way using `npm run serve`.

The frontend should use different URLs to access the server endpoints, depending on how
it has been launched.
* When launched in development mode, with `npm run serve`, it should access a go server
which is running in parallel, usually on `http://localhost:6969/`.
* In all other cases, it should assume that it has been launched from the go server,
so it should dynamically use its own URL also for calling endpoints.
This is achieved by defining a file `.env.development`.
  
## Build

To build the frontend (vue.js) part for production, issue
```npm run build```

To build the go server afterwards, issue
```docker build -t docker-hosted.financial.com/bb/grpcox:latest .```

To run the docker image, issue
```docker run -p 6969:6969 docker-hosted.financial.com/bb/grpcox:latest```

Please note that the image name `docker-hosted.financial.com/bb/grpcox` contains the
name of financial.com internal docker registry host. The commands above do not actually
contact that server, so they should work for you even if you do not have access to
financial.com internal network. However, `docker push` will certainly fail in that case,
so you might want to change the host and image name to values which suit your needs.
