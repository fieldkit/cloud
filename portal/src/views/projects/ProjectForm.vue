<template>
    <div class="project-form">
        <div class="header-row">
            <h2 v-if="!project">New Project</h2>
            <h2 v-if="project && project.id">Edit Project</h2>

            <div class="close-form-button" v-on:click="closeForm">
                <img alt="Close" src="@/assets/close.png" />
            </div>
        </div>

        <form id="form" @submit.prevent="saveForm">
            <div class="outer-input-container">
                <TextField v-model="form.name" label="Project Name" />

                <div class="validation-errors" v-if="$v.form.name.$error">
                    <div v-if="!$v.form.name.required">Name is a required field.</div>
                </div>
            </div>
            <div class="outer-input-container">
                <TextField v-model="form.description" label="Short Description" />

                <div class="validation-errors" v-if="$v.form.description.$error">
                    <div v-if="!$v.form.description.required">This is a required field.</div>
                </div>
            </div>
            <div class="outer-input-container">
                <TextField v-model="form.goal" label="Project Goal" />

                <div class="validation-errors" v-if="$v.form.goal.$error">
                    <div v-if="!$v.form.goal.required">Project goal is a required field.</div>
                </div>
            </div>
            <div class="image-container">
                <ImageUploader :image="{ url: project ? project.photo : null }" @change="onImage" />
            </div>
            <div class="outer-input-container">
                <TextField v-model="form.location" label="Location" />

                <div class="validation-errors" v-if="$v.form.location.$error">
                    <div v-if="!$v.form.location.required">Location is a required field.</div>
                </div>
            </div>

            <div class="dates-row">
                <div class="outer-date-container">
                    <div class="date-container">
                        <div class="outer-input-container">
                            <TextField v-model="form.startTime" label="Start" />
                        </div>
                        <v-date-picker
                            :value="form.pickedStart"
                            @input="updateStart"
                            :popover="{ placement: 'bottom', visibility: 'click' }"
                        >
                            <button type="button">
                                <img alt="Calendar" src="@/assets/calendar.png" />
                            </button>
                        </v-date-picker>
                    </div>

                    <div class="validation-errors" v-if="$v.form.startTime.$error">
                        <div v-if="!$v.form.startTime.date">Please enter a valid date.</div>
                    </div>
                </div>
                <div class="space-hack"></div>
                <div class="outer-date-container">
                    <div class="date-container">
                        <div class="outer-input-container">
                            <TextField v-model="form.endTime" label="End" />
                        </div>
                        <v-date-picker :value="form.pickedEnd" @input="updateEnd" :popover="{ placement: 'bottom', visibility: 'click' }">
                            <button type="button">
                                <img alt="Calendar" src="@/assets/calendar.png" />
                            </button>
                        </v-date-picker>
                    </div>

                    <div class="validation-errors" v-if="$v.form.endTime.$error">
                        <div v-if="!$v.form.endTime.date">Please enter a valid date.</div>
                    </div>
                </div>
            </div>

            <div class="outer-input-container">
                <vue-tags-input v-model="form.tag" :tags="form.tags" @tags-changed="onTagsChanged" />

                <div class="validation-errors" v-if="$v.form.tags.$error"></div>
            </div>
            <div id="public-checkbox-container">
                <input type="checkbox" id="checkbox" v-model="form.publicProject" />
                <label for="checkbox">Make this project public</label>
            </div>
            <div class="action-container">
                <button class="save-button" v-if="!project" type="submit">Add</button>
                <button class="save-button" v-if="project && project.id" type="submit">Update</button>
                <div v-if="project && project.id" class="delete-container" v-on:click="deleteProject">
                    <img alt="Delete" src="@/assets/Delete.png" />
                    Delete this project
                </div>
            </div>
        </form>
    </div>
</template>

<script lang="ts">
import _ from "lodash";
import moment from "moment";
import Vue from "@/store/strong-vue";
import StandardLayout from "../StandardLayout.vue";
import CommonComponents from "@/views/shared";
import VueTagsInput from "@johmun/vue-tags-input";

import { tryParseTags } from "@/utilities";

import { helpers, required, email, minValue, minLength, sameAs } from "vuelidate/lib/validators";

import FKApi from "@/api/api";
import * as ActionTypes from "@/store/actions";

export default Vue.extend({
    name: "ProjectForm",
    components: {
        VueTagsInput,
        ...CommonComponents,
    },
    props: {
        project: {
            type: Object,
            required: false,
        },
    },
    data: () => {
        return {
            image: null,
            form: {
                name: "",
                description: "",
                goal: "",
                location: "",
                startTime: "",
                endTime: "",
                tags: [],
                tag: "",
                publicProject: false,
                pickedStart: null,
                pickedEnd: null,
            },
        };
    },
    validations: {
        form: {
            name: {
                required,
            },
            description: {
                required,
            },
            goal: {
                required,
            },
            location: {
                required,
            },
            startTime: {
                date: function (value) {
                    if (value && value.length > 0) {
                        return moment(value).isValid();
                    }
                    return true;
                },
            },
            endTime: {
                date: function (value) {
                    if (value && value.length > 0) {
                        return moment(value).isValid();
                    }
                    return true;
                },
            },
            tags: {},
        },
    },
    mounted(this: any) {
        if (this.project) {
            this.form = {
                name: this.project.name,
                description: this.project.description,
                goal: this.project.goal,
                location: this.project.location,
                startTime: this.prettyDate(this.project.startTime),
                endTime: this.prettyDate(this.project.endTime),
                tags: tryParseTags(this.project.tags),
                publicProject: !this.project.private,
            };
        }
    },
    computed: {
        imageUrl(this: any) {
            if (this.project.photo) {
                return this.$config.baseUrl + this.project.photo;
            }
            return null;
        },
    },
    methods: {
        saveForm() {
            this.$v.form.$touch();
            if (this.$v.form.$pending || this.$v.form.$error) {
                console.log("save form, validation error");
                return;
            }

            if (this.project && this.project.id) {
                return this.updateProject();
            }
            return this.addProject();
        },
        onTagsChanged(newTags) {
            this.form.tags = newTags;
        },
        createParams() {
            const makeLocalTime = (str) => {
                if (!str || str.length == 0) {
                    return null;
                }
                return moment(str, "M/D/YYYY").toISOString();
            };
            return _.extend({}, this.form, {
                id: this.project?.id || null,
                private: !this.form.publicProject,
                startTime: makeLocalTime(this.form.startTime),
                endTime: makeLocalTime(this.form.endTime),
                tags: JSON.stringify(this.form.tags.map((tag) => tag.text)),
            });
        },
        addProject() {
            this.$emit("updating");

            const data = this.createParams();
            if (this.image) {
                return new FKApi().addProject(data).then((project) => {
                    const params = {
                        type: this.image.type,
                        file: this.image.file,
                        id: project.id,
                    };
                    return new FKApi().uploadProjectImage(params).then(() => {
                        return this.$router.push({
                            name: "viewProject",
                            params: { id: project.id },
                        });
                    });
                });
            } else {
                return new FKApi().addProject(data).then((project) => {
                    return this.$router.push({
                        name: "viewProject",
                        params: { id: project.id },
                    });
                });
            }
        },
        updateProject() {
            this.$emit("updating");

            const data = this.createParams();
            if (this.image) {
                const payload = {
                    type: this.image.type,
                    file: this.image.file,
                    id: this.project.id,
                };
                return new FKApi().uploadProjectImage(payload).then(() => {
                    return this.$store.dispatch(ActionTypes.SAVE_PROJECT, data).then(() => {
                        return this.$router.push({
                            name: "viewProject",
                            params: { id: this.project.id },
                        });
                    });
                });
            } else {
                return this.$store.dispatch(ActionTypes.SAVE_PROJECT, data).then(() => {
                    return this.$router.push({
                        name: "viewProject",
                        params: { id: this.project.id },
                    });
                });
            }
        },
        deleteProject() {
            if (window.confirm("Are you sure you want to delete this project?")) {
                return new FKApi().deleteProject({ projectId: this.project.id }).then(() => {
                    return this.$router.push({ name: "projects" });
                });
            }
        },
        updateStart(date, ...args) {
            this.form.startTime = date ? moment(date).format("M/D/YYYY") : "";
        },
        updateEnd(date, ...args) {
            this.form.endTime = date ? moment(date).format("M/D/YYYY") : "";
        },
        prettyDate(date) {
            if (date) {
                return moment(date).format("M/D/YYYY");
            }
            return "";
        },
        closeForm() {
            if (this.project) {
                return this.$router.push({ name: "viewProject", params: { id: this.project.id } });
            } else {
                return this.$router.push({ name: "projects" });
            }
        },
        onImage(image) {
            this.image = image;
        },
    },
});
</script>

<style scoped>
.project-form {
    max-width: 800px;
    display: flex;
    flex-direction: column;
    border: 1px solid #d8dce0;
    background: white;
    padding: 20px;
}

form > div {
    margin-bottom: 20px;
}

.header-row {
    display: flex;
    flex-direction: row;
}

.dates-row {
    display: flex;
    flex-direction: row;
}
.dates-row > div {
    flex: 1;
}
.image-container {
    width: 100%;
    margin: 15px 0;
}
/deep/ .image-container img {
    max-width: 275px;
    max-height: 135px;
    margin-right: 10px;
}
.date-container {
    flex: 1;
    display: flex;
}
.date-container .outer-input-container {
    flex-grow: 1;
}
.date-container button {
    position: absolute;
    margin: -4px 0 0 -30px;
    background: none;
    padding: 0;
    border: none;
}
.date-container img {
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

.outer-date-container {
    display: flex;
    flex-direction: column;
}

.action-container {
    display: flex;
    align-items: baseline;
}
.save-button {
    width: 300px;
    height: 50px;
    font-size: 18px;
    color: white;
    background-color: #ce596b;
    border: none;
    border-radius: 5px;
}
.close-form-button {
    margin-left: auto;
    margin-top: 15px;
    cursor: pointer;
}
.delete-container {
    margin-left: auto;
    cursor: pointer;
    display: inline-block;
}
.delete-container img {
    width: 12px;
    margin-right: 4px;
}
.validation-errors {
    color: #c42c44;
    display: block;
    font-size: 14px;
    margin-bottom: 25px;
}

/deep/ .ti-tag {
    background-color: #0a67aa;
}
</style>
