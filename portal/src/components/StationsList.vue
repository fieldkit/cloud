<template>
    <div id="stations-list-container">
        <div v-for="station in stations" :key="station.id" class="station-container">
            <div class="left">
                <img alt="Station image"
                    src="../assets/placeholder_station_thumbnail.png"
                    class="station-element" />
            </div>
            <div class="right">
                <div id="station-name" class="station-element">{{ station.name }}</div>
                <div class="station-element">Last Synced {{ station.synced }}</div>
                <div class="station-element">
                    <img id="battery" alt="Battery level" :src="getBatteryImg(station)" />
                    <span>{{ station.status_json.battery_level }}%</span>
                </div>
                <div class="">
                    <img v-for="module in station.status_json.moduleObjects"
                        alt="Module icon"
                        class="small-space"
                        :src="getModuleImg(module)" />
                </div>
            </div>
        </div>
    </div>
</template>

<script>
export default {
    name: "StationsList",
    data: () => {
        return {
        };
    },
    props: ['stations'],
    mounted(args) {
        this.stations.forEach((s) => {
            this.$set(s, 'synced', this.formatDate(s.status_json.updated));
        });
    },
    methods: {
        formatDate(date) {
            let d = new Date(date);
            return d.toLocaleDateString("en-US");
        },
        getBatteryImg(station) {
            let images = require.context('../assets/battery/', false, /\.png$/);
            let battery = station.status_json.battery_level;
            // check to see if it already has a percent sign
            if (battery.toString().indexOf("%") > -1) {
                battery = parseInt(
                    battery.toString().split("%")[0]
                );
            }
            let img = "";
            if(battery == 0) {
                img = "0.png";
            } else if(battery <= 20) {
                img = "20.png";
            } else if(battery <= 40) {
                img = "40.png";
            } else if(battery <= 60) {
                img = "60.png";
            } else if(battery <= 80) {
                img = "80.png";
            } else {
                img = "100.png";
            }
            return images('./' + img)
        },
        getModuleImg(module) {
            let images = require.context('../assets/', false, /\.png$/);
            let img = "placeholder.png";
            // Note: this is not a trustworthy way of figuring out what icons to show,
            // as the user could rename their module anything
            if(module.name.indexOf("Water") > -1) {
                img = "water.png";
            }
            if(module.name.indexOf("Weather") > -1) {
                img = "weather.png";
            }
            if(module.name.indexOf("Ocean") > -1) {
                img = "ocean.png";
            }
            return images('./' + img);
        }
    }
};
</script>

<style scoped>
#stations-list-container {
    width: 700px;
    margin: 60px 0;
}
.station-container {
    width: 400px;
    padding: 10px;
    margin: 20px 0;
    border: 1px solid rgb(215, 220, 225);
    font-size: 14px;
    font-weight: lighter;
    overflow: hidden;
}
.left, .right {
    float: left;
}
.station-element {
    margin: 5px 5px 0 5px;
}
#station-name {
    font-size: 18px;
    font-weight: bold;
}
#battery {
    margin-right: 5px;
}
.small-space {
    margin: 3px;
}
</style>
