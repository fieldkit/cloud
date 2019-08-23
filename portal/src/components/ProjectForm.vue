<template>
    <div id="project-form-container">
        <div id="close-form-btn" v-on:click="closeForm">
            <img alt="Close" src="../assets/close.png" />
        </div>
        <h2>New Project</h2>
        <input v-model="name" placeholder="Project Name" class="text-input wide-text-input" />
        <input v-model="description" placeholder="Short Description" class="text-input wide-text-input" />
        <input v-model="goal" placeholder="Project Goal" class="text-input wide-text-input" />
        <div id="add-image">
            <img alt="Add project image" src="../assets/add_image.png" />
        </div>
        <input v-model="location" placeholder="Location" class="text-input wide-text-input" />
        <div id="start-date">
            <input v-model="startDate" placeholder="Start Date" class="text-input" />
            <v-date-picker v-model="startDate" :popover="{ placement: 'bottom', visibility: 'click' }">
                <button>
                    <img alt="Calendar" src="../assets/calendar.png" />
                </button>
            </v-date-picker>
        </div>
        <div id="end-date">
            <input v-model="endDate" placeholder="End Date" class="text-input" />
            <v-date-picker v-model="endDate" :popover="{ placement: 'bottom', visibility: 'click' }">
                <button>
                    <img alt="Calendar" src="../assets/calendar.png" />
                </button>
            </v-date-picker>
        </div>
        <input v-model="tags" placeholder="Add Tags" class="text-input wide-text-input" />
        <div id="public-checkbox-container">
            <input type="checkbox" id="checkbox" v-model="publicProject" />
            <label for="checkbox">Make this project public</label>
            <img alt="Info" src="../assets/info.png" />
        </div>
        <button class="save-btn" :disabled="isDisabled">Save</button>
    </div>
</template>

<script>
export default {
    name: "ProjectForm",
    data: () => {
        return {
            name: "",
            description: "",
            goal: "",
            imageName: "",
            location: "",
            startDate: "",
            endDate: "",
            tags: "",
            publicProject: false,
            isDisabled: true
        };
    },
    methods: {
        closeForm() {
            this.$emit("closeProjectForm");
        }
    }
};
</script>

<style scoped>
#project-form-container {
    width: 700px;
    padding: 0 15px 15px 15px;
    margin: 60px 0;
    border: 1px solid rgb(215, 220, 225);
}
.wide-text-input {
    width: 98%;
    border: none;
    border-bottom: 2px solid rgb(235, 235, 235);
    font-size: 15px;
    margin: 15px 0;
    padding-bottom: 4px;
}
#add-image {
    margin: 15px 0;
}
#start-date input,
#end-date input {
    width: 90%;
    border: none;
    font-size: 15px;
    margin: 0;
    padding: 0;
    vertical-align: bottom;
}
#start-date button,
#end-date button {
    float: right;
    background: none;
    padding: 0;
    margin: 0;
}
#start-date img,
#end-date img {
    vertical-align: bottom;
    padding-bottom: 2px;
    margin-left: 4px;
}
#start-date,
#end-date {
    line-height: 1.75;
    width: 47%;
    margin: 15px 0;
    border-bottom: 2px solid rgb(235, 235, 235);
}
#start-date {
    float: left;
}
#end-date {
    float: right;
}
#public-checkbox-container {
    margin: 15px 0;
    width: 98%;
    padding-bottom: 20px;
}
#public-checkbox-container input {
    float: left;
    margin: 5px;
}
#public-checkbox-container label {
    float: left;
    margin: 2px 5px;
}
#public-checkbox-container img {
    float: left;
    margin: 2px 5px;
}

.save-btn {
    width: 300px;
    height: 45px;
    font-size: 18px;
    font-weight: bold;
    color: white;
    background: rgb(215, 220, 225);
    border: none;
    border-radius: 6px;
    margin: 60px 0 20px 0;
}

#close-form-btn {
    float: right;
    margin-top: 15px;
    cursor: pointer;
}
</style>
