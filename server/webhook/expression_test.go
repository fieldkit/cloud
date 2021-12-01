package webhook

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/elgs/gojq"
)

var Payload1 = `[
        {
            "battery": 4.288,
            "distance": 2715,
            "sdError": 0,
            "solarboard_voltage": 6.6,
            "counter": 186415,
            "port": 1,
            "frequency": 903.9,
            "airtime": 185344000,
            "dt": "2021-09-24T19:36:04.802813185Z"
        },
        {
            "app_id": "deployment_two_app",
            "dev_id": "sensor_8",
            "hw_serial": "002507D3349747CE",
            "data_rate": "SF9BW125"
        }
    ]`

var Payload2 = `[
        {
            "dt": "2021-09-24T19:35:12.112654658Z",
            "rssi": -109,
            "snr": -5.25,
            "bw": 125000,
            "sf": 10,
            "freq": "905100000",
            "port": 1,
            "count": 23211,
            "altitude": 0,
            "battery": 4.163,
            "distance": 3042,
            "humidity": 0,
            "pressure": 0,
            "sdError": 0,
            "temperature": 0,
            "depth": 0
        },
        {
            "dev_id": "smith-and-9th",
            "app_id": "deployment-3",
            "dev_addr": "260C6160",
            "dev_eui": "AADCCCCDC4343445"
        }
    ]`

var Payload3 = `{
        "end_device_ids": {
            "device_id": "smith-and-9th",
            "application_ids": {
                "application_id": "deployment-3"
            },
            "dev_eui": "EFDC465458898445",
            "join_eui": "0000000000000000",
            "dev_addr": "25646589"
        },
        "received_at": "2021-09-24T19:35:12.402737172Z",
        "uplink_message": {
            "f_port": 1,
            "f_cnt": 23211,
            "frm_payload": "AEMQ4gs=",
            "decoded_payload": {
                "altitude": 0,
                "battery": 4.163,
                "distance": 3042,
                "humidity": 0,
                "pressure": 0,
                "sdError": 0,
                "temperature": 0
            },
            "settings": {
                "data_rate": {
                    "lora": {
                        "bandwidth": 125000,
                        "spreading_factor": 10
                    }
                },
                "coding_rate": "4/5",
                "frequency": "905100000",
                "timestamp": 220402788
            },
            "received_at": "2021-09-24T19:35:12.112654658Z",
            "consumed_airtime": "0.329728s",
            "network_ids": {
                "net_id": "000013",
                "tenant_id": "ttn",
                "cluster_id": "ttn-nam1"
            }
        }
    },
]`

func TestSimpleExpression(t *testing.T) {
	assert := assert.New(t)

	parser, err := gojq.NewStringQuery(Payload1)
	assert.NoError(err)

	value, err := parser.Query("[1].hw_serial")
	assert.NoError(err)

	assert.NotNil(value)
}
