<template>
    <div id="project-form-container">
        <div id="close-form-btn" v-on:click="closeForm">
            <img alt="Close" src="../assets/close.png" />
        </div>
        <h2>{{ this.formHeading }}</h2>
        <div class="outer-input-container">
            <div class="input-container">
                <input v-model="name" class="inputText" required="" />
                <span class="floating-label">Project Name</span>
            </div>
        </div>
        <div class="outer-input-container">
            <div class="input-container">
                <input v-model="description" class="inputText" required="" />
                <span class="floating-label">Short Description</span>
            </div>
        </div>
        <div class="outer-input-container">
            <div class="input-container">
                <input v-model="goal" class="inputText" required="" />
                <span class="floating-label">Project Goal</span>
            </div>
        </div>
        <div class="image-container">
            <img alt="Project image" :src="imageUrl" class="custom-image" v-if="hasImage && !previewImage" />
            <img :src="previewImage" class="custom-image" v-if="!hasImage || previewImage" />
            <img alt="Add project image" src="../assets/add_image.png" class="custom-image" v-if="!hasImage && !previewImage" />
            <div class="image-instruct">{{ this.hasImage ? "Change image " : "Add an image to your project " }}</div>
            <input type="file" accept="image/gif, image/jpeg, image/png" @change="uploadImage" />
        </div>
        <div class="outer-input-container">
            <div class="input-container">
                <input v-model="location" class="inputText" required="" />
                <span class="floating-label">Location</span>
            </div>
        </div>
        <div id="start-date">
            <div class="outer-input-container">
                <div class="input-container">
                    <input v-model="displayStartDate" class="inputText" required="" />
                    <span class="floating-label">Start Date</span>
                </div>
            </div>
            <v-date-picker v-model="startDate" @input="updateDisplayDates" :popover="{ placement: 'bottom', visibility: 'click' }">
                <button>
                    <img alt="Calendar" src="../assets/calendar.png" />
                </button>
            </v-date-picker>
        </div>
        <div id="end-date">
            <div class="outer-input-container">
                <div class="input-container">
                    <input v-model="displayEndDate" class="inputText" required="" />
                    <span class="floating-label">End Date</span>
                </div>
            </div>
            <v-date-picker v-model="endDate" @input="updateDisplayDates" :popover="{ placement: 'bottom', visibility: 'click' }">
                <button>
                    <img alt="Calendar" src="../assets/calendar.png" />
                </button>
            </v-date-picker>
        </div>
        <div class="outer-input-container">
            <div class="input-container">
                <input v-model="tags" class="inputText" required="" />
                <span class="floating-label">Tags</span>
            </div>
        </div>
        <div id="public-checkbox-container">
            <input type="checkbox" id="checkbox" v-model="publicProject" />
            <label for="checkbox">Make this project public</label>
            <img alt="Info" src="../assets/info.png" />
        </div>
        <div class="action-container">
            <button class="save-btn" v-if="formType == 'add'" v-on:click="addProject">
                Add
            </button>
            <button class="save-btn" v-if="formType == 'update'" v-on:click="updateProject">
                Update
            </button>
            <div v-if="formType == 'update'" class="delete-container" v-on:click="deleteProject">
                <img alt="Info" src="../assets/Delete.png" />
                Delete this project
            </div>
        </div>
    </div>
</template>

<script>
import FKApi from "../api/api";
import Config from "../secrets";

export default {
    name: "ProjectForm",
    props: ["project"],
    data: () => {
        return {
            formType: "add",
            formHeading: "New Project",
            name: "",
            description: "",
            goal: "",
            imageName: "",
            location: "",
            startDate: null,
            displayStartDate: "",
            endDate: null,
            displayEndDate: "",
            tags: "",
            publicProject: false,
            hasImage: false,
            imageUrl: "",
            baseUrl: Config.API_HOST,
            previewImage: null,
            acceptedImageTypes: ["jpg", "jpeg", "png", "gif"],
        };
    },
    watch: {
        project(_project) {
            if (_project) {
                this.formHeading = "Edit Project";
                this.formType = "update";
                this.name = _project.name;
                this.description = _project.description;
                this.goal = _project.goal;
                this.location = _project.location;
                this.startDate = _project.start_time;
                this.endDate = _project.end_time;
                this.tags = _project.tags;
                this.publicProject = !_project.private;
                this.updateDisplayDates();
                if (_project.media_url) {
                    this.imageUrl = this.baseUrl + "/projects/" + _project.id + "/media";
                    this.hasImage = true;
                } else {
                    this.imageUrl = "";
                    this.hasImage = false;
                }
            } else {
                this.formType = "add";
                this.formHeading = "New Project";
                this.resetFields();
            }
        },
    },
    methods: {
        createParams() {
            const data = {
                description: this.description,
                end_time: this.endDate,
                goal: this.goal,
                location: this.location,
                name: this.name,
                private: !this.publicProject,
                slug: "proj-" + Date.now(),
                start_time: this.startDate,
                tags: this.tags,
            };
            if (this.project) {
                data.id = this.project.id;
                data.slug = this.project.slug;
            }
            return data;
        },
        addProject() {
            this.$emit("updating");
            const api = new FKApi();
            const data = this.createParams();
            if (this.sendingImage) {
                api.addProject(data).then(project => {
                    let params = {
                        type: this.imageType,
                        image: this.sendingImage,
                        id: project.id,
                    };
                    api.uploadProjectImage(params).then(() => {
                        this.$router.push({
                            name: "viewProject",
                            params: { id: project.id },
                        });
                    });
                });
            } else {
                api.addProject(data).then(project => {
                    this.$router.push({
                        name: "viewProject",
                        params: { id: project.id },
                    });
                });
            }
        },
        updateProject() {
            this.$emit("updating");
            const api = new FKApi();
            const data = this.createParams();
            if (this.sendingImage) {
                let params = {
                    type: this.imageType,
                    image: this.sendingImage,
                    id: this.project.id,
                };
                api.uploadProjectImage(params).then(() => {
                    api.updateProject(data).then(project => {
                        this.$router.push({
                            name: "viewProject",
                            params: { id: project.id },
                        });
                    });
                });
            } else {
                api.updateProject(data).then(project => {
                    this.$router.push({
                        name: "viewProject",
                        params: { id: project.id },
                    });
                });
            }
        },
        deleteProject() {
            if (window.confirm("Are you sure you want to delete this project?")) {
                const api = new FKApi();
                const params = {
                    projectId: this.project.id,
                };
                api.deleteProject(params).then(() => {
                    this.$router.push({ name: "projects" });
                });
            }
        },
        updateDisplayDates() {
            if (this.startDate) {
                let d = new Date(this.startDate);
                this.displayStartDate = d.toLocaleDateString("en-US");
            }
            if (this.endDate) {
                let d = new Date(this.endDate);
                this.displayEndDate = d.toLocaleDateString("en-US");
            }
        },
        resetFields() {
            this.name = "";
            this.description = "";
            this.goal = "";
            this.imageName = "";
            this.location = "";
            this.startDate = null;
            this.displayStartDate = "";
            this.endDate = null;
            this.displayEndDate = "";
            this.tags = "";
            this.publicProject = false;
            this.hasImage = false;
            this.imageUrl = "";
            this.previewImage = null;
        },
        uploadImage(event) {
            this.previewImage = null;
            this.sendingImage = null;
            let valid = false;
            if (event.target.files.length > 0) {
                this.acceptedImageTypes.forEach(t => {
                    if (event.target.files[0].type.indexOf(t) > -1) {
                        valid = true;
                    }
                });
            }
            if (!valid) {
                return;
            }
            this.imageType = event.target.files[0].type;
            const image = event.target.files[0];
            this.sendingImage = image;
            const reader = new FileReader();
            reader.readAsDataURL(image);
            reader.onload = event => {
                this.previewImage = event.target.result;
            };
        },
        closeForm() {
            this.$router.push({ name: "projects" });
        },
    },
};
</script>

<style scoped>
#project-form-container {
    width: 700px;
    float: left;
    padding: 0 15px 15px 15px;
    margin: 25px 0;
    border: 1px solid rgb(215, 220, 225);
}
.outer-input-container {
    float: left;
    width: 98%;
    margin: 19px 0 0 0;
}
#start-date .outer-input-container,
#end-date .outer-input-container {
    width: 100%;
    margin: 0;
}
.input-container {
    margin: auto;
    width: 100%;
    text-align: left;
}
.image-container {
    width: 98%;
    margin: 15px 0;
    float: left;
}
.custom-image {
    max-width: 275px;
    max-height: 135px;
    margin-right: 10px;
}
.image-instruct {
    font-size: 14px;
    margin-right: 5px;
    display: inline-block;
}

#start-date,
#end-date {
    width: 47%;
    margin: 19px 0 0 0;
}
#start-date {
    float: left;
}
#end-date {
    float: right;
}
#start-date button,
#end-date button {
    position: absolute;
    margin: -4px 0 0 -30px;
    background: none;
    padding: 0;
    border: none;
}
#start-date img,
#end-date img {
    vertical-align: bottom;
    padding-bottom: 2px;
    margin-left: 4px;
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

.action-container {
    float: left;
    clear: both;
}
.save-btn {
    width: 300px;
    height: 50px;
    font-size: 18px;
    color: white;
    background-color: #ce596b;
    border: none;
    border-radius: 5px;
    margin: 50px 0 20px 0;
}

#close-form-btn {
    float: right;
    margin-top: 15px;
    cursor: pointer;
}

.delete-container {
    cursor: pointer;
    margin-left: 230px;
    display: inline-block;
}
.delete-container img {
    width: 12px;
    margin-right: 4px;
}
</style>
