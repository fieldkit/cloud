module github.com/fieldkit/cloud/server

go 1.12

require (
	github.com/Nerzal/gocloak v1.0.0 // indirect
	github.com/Nerzal/gocloak/v7 v7.11.0
	github.com/O-C-R/singlepage v0.0.0-20170327184421-bbbe2159ecec
	github.com/PyYoshi/goa-logging-zap v0.2.3
	github.com/ajg/form v1.5.1
	github.com/araddon/dateparse v0.0.0-20200409225146-d820a6159ab1
	github.com/armon/go-metrics v0.0.0-20180917152333-f0300d1749da
	github.com/aws/aws-lambda-go v0.0.0-20180413184133-ea03c2814414
	github.com/aws/aws-sdk-go v1.29.20
	github.com/bgentry/que-go v1.0.1
	github.com/bxcodec/faker/v3 v3.3.1
	github.com/cavaliercoder/grab v2.0.0+incompatible
	github.com/cenkalti/backoff v0.0.0-20170309153948-3db60c813733
	github.com/conservify/gonaturalist v0.0.0-20190530183130-1509fd074b2c
	github.com/conservify/protobuf-tools v0.0.0-20180715164506-43b897198d14
	github.com/conservify/tooling v0.0.0-20190530175209-bf2b6e69e188
	github.com/coreos/go-oidc v2.2.1+incompatible
	github.com/creack/goselect v0.1.0 // indirect
	github.com/crewjam/saml v0.4.2
	github.com/davecgh/go-spew v1.1.1
	github.com/dghubble/go-twitter v0.0.0-20170219215544-fc93bb35701b
	github.com/dghubble/oauth1 v0.0.0-20170219195226-3c7784d12fed
	github.com/dghubble/sling v0.0.0-20170219194632-91b015f8c5e2
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/dimfeld/httppath v0.0.0-20170720192232-ee938bf73598
	github.com/dimfeld/httptreemux v5.0.1+incompatible
	github.com/disintegration/gift v1.2.1 // indirect
	github.com/disintegration/imageorient v0.0.0-20180920195336-8147d86e83ec
	github.com/dropbox/godropbox v0.0.0-20190501155911-5749d3b71cbe // indirect
	github.com/elgs/gojq v0.0.0-20201120033525-b5293fef2759
	github.com/elgs/gosplitargs v0.0.0-20161028071935-a491c5eeb3c8 // indirect
	github.com/fieldkit/app-protocol v0.0.0-20200515173549-e0925480073d
	github.com/fieldkit/data-protocol v0.0.0-20200825203301-d452c6499956
	github.com/go-ini/ini v1.26.0
	github.com/go-kit/kit v0.8.0
	github.com/go-logfmt/logfmt v0.5.0
	github.com/go-pg/pg/v10 v10.9.1
	github.com/go-stack/stack v1.8.0
	github.com/goadesign/goa v1.4.0
	github.com/gogo/protobuf v1.3.1
	github.com/golang-migrate/migrate/v4 v4.11.0
	github.com/golang/freetype v0.0.0-20161208064710-d9be45aaf745
	github.com/golang/protobuf v1.4.3
	github.com/google/go-querystring v1.0.0
	github.com/google/uuid v1.1.2
	github.com/gorilla/mux v1.8.0
	github.com/gorilla/websocket v1.4.2
	github.com/govau/que-go v1.0.1
	github.com/h2non/filetype v1.0.10
	github.com/hashicorp/go-immutable-radix v1.0.0
	github.com/hashicorp/go-multierror v1.1.0
	github.com/hashicorp/golang-lru v0.5.1
	github.com/hongshibao/go-algo v0.0.0-20160521171829-b1aaa26798b6 // indirect
	github.com/hongshibao/go-kdtree v0.0.0-20180503061502-0de4e8305acf
	github.com/iancoleman/strcase v0.0.0-20190422225806-e506e3ef7365
	github.com/inconshreveable/log15 v0.0.0-20180818164646-67afb5ed74ec
	github.com/inconshreveable/mousetrap v1.0.0
	github.com/influxdata/influxdb-client-go/v2 v2.8.1 // indirect
	github.com/itchyny/gojq v0.12.5
	github.com/jackc/pgx v3.6.2+incompatible
	github.com/jackc/pgx/v4 v4.16.1 // indirect
	github.com/jmespath/go-jmespath v0.0.0-20180206201540-c2b33e8439af
	github.com/jmoiron/sqlx v1.2.0
	github.com/kelseyhightower/envconfig v0.0.0-20170206223400-8bf4bbfc795e
	github.com/kinbiko/jsonassert v1.0.1
	github.com/konsorten/go-windows-terminal-sequences v1.0.2
	github.com/kr/logfmt v0.0.0-20140226030751-b84e30acd515
	github.com/lib/pq v1.10.4
	github.com/llgcode/draw2d v0.0.0-20161104081029-1286d3b2030a
	github.com/lucasb-eyer/go-colorful v0.0.0-20170223221042-c900de9dbbc7
	github.com/manveru/faker v0.0.0-20171103152722-9fbc68a78c4d
	github.com/matoous/go-nanoid/v2 v2.0.0 // indirect
	github.com/mattn/go-colorable v0.1.8
	github.com/mattn/go-isatty v0.0.13
	github.com/montanaflynn/stats v0.5.0
	github.com/muesli/smartcrop v0.3.0
	github.com/nfnt/resize v0.0.0-20180221191011-83c6a9932646 // indirect
	github.com/paulmach/go.geo v0.0.0-20170321183534-b160a6efed6c
	github.com/paulmach/go.geojson v0.0.0-20170327170536-40612a87147b
	github.com/paulmach/orb v0.1.6
	github.com/pkg/errors v0.9.1
	github.com/pkg/profile v1.5.0
	github.com/pquerna/cachecontrol v0.0.0-20201205024021-ac21108117ac // indirect
	github.com/robinjoseph08/go-pg-migrations/v3 v3.0.0
	github.com/robinpowered/go-proto v0.0.0-20160614142341-85ea3e1f1d3d
	github.com/satori/go.uuid v1.2.1-0.20181028125025-b2ce2384e17b
	github.com/segmentio/ksuid v1.0.3
	github.com/sirupsen/logrus v1.4.2
	github.com/sohlich/go-dbscan v0.0.0-20161128164835-242a0c72bf77
	github.com/spf13/cobra v0.0.3
	github.com/spf13/pflag v1.0.3
	github.com/spf13/viper v1.7.1
	github.com/stretchr/testify v1.7.0
	github.com/tidwall/gjson v1.6.0
	github.com/ugorji/go v0.0.0-20190320090025-2dc34c0b8780
	github.com/xeipuuv/gojsonpointer v0.0.0-20170225233418-6fe8760cad35
	github.com/xeipuuv/gojsonreference v0.0.0-20150808065054-e02fc20de94c
	github.com/xeipuuv/gojsonschema v0.0.0-20170324002221-702b404897d4
	github.com/zach-klippenstein/goregen v0.0.0-20160303162051-795b5e3961ea
	go.bug.st/serial.v1 v0.0.0-20180827123349-5f7892a7bb45 // indirect
	go.uber.org/atomic v1.6.0
	go.uber.org/multierr v1.5.0
	go.uber.org/zap v1.13.0
	goa.design/goa/v3 v3.2.4
	goa.design/plugins/v3 v3.1.1
	golang.org/x/crypto v0.0.0-20210711020723-a769d52b0f97
	golang.org/x/image v0.0.0-20190802002840-cff245a6509b
	golang.org/x/net v0.0.0-20210226172049-e18ecbb05110
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d
	golang.org/x/sys v0.0.0-20210831042530-f4d43177bf5e
	golang.org/x/tools v0.0.0-20201009162240-fcf82128ed91
	google.golang.org/appengine v1.6.7
	google.golang.org/genproto v0.0.0-20201009135657-4d944d34d83c // indirect
	google.golang.org/grpc v1.33.0 // indirect
	gopkg.in/alexcesaro/statsd.v2 v2.0.0
	gopkg.in/square/go-jose.v2 v2.5.1 // indirect
	gopkg.in/yaml.v2 v2.3.0
)
