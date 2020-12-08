<template>
    <div class="form-edit">
        <div class="header-row">
            <h2 v-if="!project">New Project</h2>
            <h2 v-if="project && project.id">Edit Project</h2>

            <div class="close-form-button" v-on:click="closeForm">
                <img alt="Close" src="@/assets/icon-close.svg" />
            </div>
        </div>

        <form id="form" @submit.prevent="saveForm">
            <div class="outer-input-container">
                <TextField v-model="form.name" label="Project Name" />

                <div class="validation-errors" v-if="$v.form.name.$error">
                    <div v-if="!$v.form.name.required">Name is a required field.</div>
                    <div v-if="!$v.form.name.maxLength">This field has a limit of 100 characters.</div>
                </div>
            </div>
            <div class="outer-input-container">
                <TextField v-model="form.description" label="Short Description" />

                <div class="validation-errors" v-if="$v.form.description.$error">
                    <div v-if="!$v.form.description.required">This is a required field.</div>
                    <div v-if="!$v.form.description.maxLength">This field has a limit of 100 characters.</div>
                </div>
            </div>
            <div class="outer-input-container">
                <TextField v-model="form.goal" label="Project Goal" />

                <div class="validation-errors" v-if="$v.form.goal.$error">
                    <div v-if="!$v.form.goal.required">Project goal is a required field.</div>
                    <div v-if="!$v.form.goal.maxLength">This field has a limit of 100 characters.</div>
                </div>
            </div>
            <div class="image-container">
                <ImageUploader :image="{ url: project ? project.photo : null }" :placeholder="imagePlaceholder" @change="onImage" />
            </div>
            <div class="outer-input-container">
                <TextField v-model="form.location" label="Location" />

                <div class="validation-errors" v-if="$v.form.location.$error">
                    <div v-if="!$v.form.location.required">Location is a required field.</div>
                    <div v-if="!$v.form.location.maxLength">This field has a limit of 100 characters.</div>
                </div>
            </div>

            <div class="dates-row">
                <div class="date-container">
                    <div class="outer-input-container">
                        <TextField v-model="form.startTime" label="Start" />
                    </div>
                    <v-date-picker :value="form.pickedStart" @input="updateStart" :popover="{ placement: 'auto', visibility: 'click' }">
                        <button type="button">
                            <img alt="Calendar" src="@/assets/icon-calendar-gray.svg" />
                        </button>
                    </v-date-picker>
                </div>

                <div class="validation-errors" v-if="$v.form.startTime.$error">
                    <div v-if="!$v.form.startTime.date">Please enter a valid date.</div>
                </div>

                <div class="date-container">
                    <div class="outer-input-container">
                        <TextField v-model="form.endTime" label="End" />
                    </div>
                    <v-date-picker :value="form.pickedEnd" @input="updateEnd" :popover="{ placement: 'auto', visibility: 'click' }">
                        <button type="button">
                            <img alt="Calendar" src="@/assets/icon-calendar-gray.svg" />
                        </button>
                    </v-date-picker>
                </div>

                <div class="validation-errors" v-if="$v.form.endTime.$error">
                    <div v-if="!$v.form.endTime.date">Please enter a valid date.</div>
                    <div v-if="!$v.form.endTime.minValue">Please enter a date after the start date.</div>
                </div>
            </div>

            <div class="outer-input-container tags-container">
                <span v-bind:class="{ focused: smallTagsLabel }">Tags</span>
                <vue-tags-input
                    v-model="form.tag"
                    :tags="form.tags"
                    @tags-changed="onTagsChanged"
                    @blur="onTagsBlur"
                    @focus="onTagsFocus"
                    placeholder=""
                />

                <div class="validation-errors" v-if="$v.form.tags.$error">
                    <div v-if="!$v.form.tags.maxLength">This field has a limit of 100 characters.</div>
                </div>
            </div>
            <div class="privacy">
                <div class="checkbox">
                    <label>
                        Make this project public
                        <input type="checkbox" id="checkbox" v-model="form.public" />
                        <span class="checkbox-btn"></span>
                    </label>
                </div>

                <div v-if="form.public" class="radio-container">
                    <label class="radio">
                        <input type="radio" id="privacy" value="1" v-model="form.privacy" />
                        <span class="radio-btn"></span>
                        Show exact location of stations
                    </label>
                    <label class="radio">
                        <input type="radio" id="privacy" value="2" v-model="form.privacy" />
                        <span class="radio-btn"></span>
                        Show general location of stations
                    </label>
                </div>
            </div>
            <div class="action-container">
                <button class="btn" v-if="!project" type="submit">Add project</button>
                <button class="btn" v-if="project && project.id" type="submit">Save updates</button>
                <button v-if="project && project.id" class="btn btn-delete" type="submit" v-on:click.prevent="deleteProject">
                    Delete Project
                </button>
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
import { UploadedImage } from "@/views/shared/ImageUploader.vue";
import { tryParseTags } from "@/utilities";

import { helpers, required, email, minValue, maxLength, minLength } from "vuelidate/lib/validators";

import * as ActionTypes from "@/store/actions";

import PlaceholderImage from "@/assets/image-placeholder.svg";

const afterOtherDate = (afterOtherDate) =>
    helpers.withParams({ type: "afterOtherDate", after: afterOtherDate }, function (this: any, value, parentVm) {
        const other = helpers.ref(afterOtherDate, this, parentVm);
        if (!other || other.length === 0) {
            return true;
        }
        if (!value || value.length == 0) {
            return true;
        }
        return moment(other).isSameOrBefore(moment(value));
    });

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
    data(): {
        image: UploadedImage | null;
        tagsFocused: boolean;
        imagePlaceholder: string;
        form: {
            name: string;
            description: string;
            goal: string;
            location: string;
            startTime: string;
            endTime: string;
            tags: { text: string }[];
            tag: string;
            public: boolean;
            privacy: number;
            pickedStart: number | null;
            pickedEnd: number | null;
        };
    } {
        return {
            image: null,
            tagsFocused: false,
            imagePlaceholder: PlaceholderImage,
            form: {
                name: "",
                description: "",
                goal: "",
                location: "",
                startTime: "",
                endTime: "",
                tags: [],
                tag: "",
                public: false,
                privacy: 1,
                pickedStart: null,
                pickedEnd: null,
            },
        };
    },
    validations: {
        form: {
            name: {
                required,
                maxLength: maxLength(100),
            },
            description: {
                required,
                maxLength: maxLength(100),
            },
            goal: {
                required,
                maxLength: maxLength(100),
            },
            location: {
                required,
                maxLength: maxLength(100),
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
                minValue: afterOtherDate("startTime"),
            },
            tags: {
                maxLength: (value) => {
                    const raw = JSON.stringify(value.map((tag) => tag.text));
                    return raw.length <= 100;
                },
            },
        },
    },
    mounted(): void {
        if (this.project) {
            this.form = {
                name: this.project.name,
                description: this.project.description,
                goal: this.project.goal,
                location: this.project.location,
                startTime: this.prettyDate(this.project.startTime),
                endTime: this.prettyDate(this.project.endTime),
                tags: tryParseTags(this.project.tags),
                public: this.project.privacy > 0,
                privacy: this.project.privacy == 0 ? 1 : this.project.privacy,
                pickedStart: null,
                pickedEnd: null,
                tag: "",
            };
        }
    },
    computed: {
        smallTagsLabel(): boolean {
            return (this.form.tags && this.form.tags.length > 0) || this.tagsFocused;
        },
        imageUrl(): string | null {
            if (this.project.photo) {
                return this.$config.baseUrl + this.project.photo;
            }
            return null;
        },
    },
    methods: {
        onTagsFocus(): void {
            this.tagsFocused = true;
        },
        onTagsBlur(): void {
            this.tagsFocused = false;
        },
        async saveForm(): Promise<void> {
            this.$v.form.$touch();
            if (this.$v.form.$pending || this.$v.form.$error) {
                console.log("save form, validation error");
                return;
            }

            console.log("saving form");

            if (this.project && this.project.id) {
                await this.updateProject();
            } else {
                await this.addProject();
            }
        },
        onTagsChanged(newTags): void {
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
                privacy: this.form.public ? Number(this.form.privacy || 0) : 0,
                startTime: makeLocalTime(this.form.startTime),
                endTime: makeLocalTime(this.form.endTime),
                tags: JSON.stringify(this.form.tags.map((tag) => tag.text)),
            });
        },
        async addProject(): Promise<void> {
            this.$emit("updating");

            const data = this.createParams();
            const image = this.image;
            if (image) {
                await this.$store.dispatch(ActionTypes.ADD_PROJECT, data).then((project) => {
                    const params = {
                        type: image.type,
                        file: image.file,
                        id: project.id,
                    };
                    return this.$services.api.uploadProjectImage(params).then(() => {
                        return this.$router.push({
                            name: "viewProject",
                            params: { id: project.id },
                        });
                    });
                });
            } else {
                await this.$store.dispatch(ActionTypes.ADD_PROJECT, data).then((project) => {
                    return this.$router.push({
                        name: "viewProject",
                        params: { id: project.id },
                    });
                });
            }
        },
        async updateProject(): Promise<void> {
            console.log("updating");

            this.$emit("updating");

            const data = this.createParams();
            const image = this.image;
            if (image) {
                const payload = {
                    type: image.type,
                    file: image.file,
                    id: this.project.id,
                };
                await this.$services.api.uploadProjectImage(payload).then(() => {
                    return this.$store.dispatch(ActionTypes.SAVE_PROJECT, data).then(() => {
                        return this.$router.push({
                            name: "viewProject",
                            params: { id: this.project.id },
                        });
                    });
                });
            } else {
                await this.$store.dispatch(ActionTypes.SAVE_PROJECT, data).then(() => {
                    return this.$router.push({
                        name: "viewProject",
                        params: { id: this.project.id },
                    });
                });
            }
        },
        async deleteProject(): Promise<void> {
            if (window.confirm("Are you sure you want to delete this project?")) {
                await this.$store.dispatch(ActionTypes.DELETE_PROJECT, { projectId: this.project.id }).then(() => {
                    return this.$router.push({ name: "projects" });
                });
            }
        },
        updateStart(date): void {
            this.form.startTime = date ? moment(date).format("M/D/YYYY") : "";
        },
        updateEnd(date, ...args): void {
            this.form.endTime = date ? moment(date).format("M/D/YYYY") : "";
        },
        prettyDate(date): string {
            if (date) {
                return moment(date).format("M/D/YYYY");
            }
            return "";
        },
        async closeForm(): Promise<void> {
            if (this.project) {
                await this.$router.push({ name: "viewProject", params: { id: this.project.id } });
            } else {
                await this.$router.push({ name: "projects" });
            }
        },
        onImage(image): void {
            this.image = image;
        },
    },
});
</script>

<style scoped lang="scss">
@import "../../scss/forms";
@import "../../scss/global";

form > .outer-input-container {
    margin-bottom: 20px;
}

.header-row {
    display: flex;
    flex-direction: row;
}

.dates-row {
    @include flex(center, space-between);
    flex-direction: row;
    margin-bottom: 20px;
}
.dates-row > div {
    flex: 1;
}
.image-container {
    width: 100%;
    margin: 28px 0 15px;
}
::v-deep .image-container .img {
    max-width: 275px;
    max-height: 135px;
    margin-right: 10px;
}
.date-container {
    flex: 1;
    display: flex;
    position: relative;
    max-width: 200px;
}
.date-container .outer-input-container {
    flex-grow: 1;
}
.date-container button {
    position: absolute;
    bottom: 5px;
    right: 0;
    margin: -4px 0 0 -30px;
    background: none;
    padding: 0;
    border: none;
}
.date-container img {
    vertical-align: bottom;
    padding-bottom: 2px;
}

.privacy {
    margin: 34px 0;
}

.radio {
    padding-left: 32px;
    margin: 7px 0;
    position: relative;
    cursor: pointer;
    min-height: 22px;
    @include flex(center);

    input {
        opacity: 0;
        position: absolute;
        z-index: -1;
    }

    &-btn {
        content: "";
        width: 20px;
        height: 20px;
        border-radius: 100px;
        border: solid 1px rgba(0, 0, 0, 0.1);
        background: #f2f4f7;
        @include position(absolute, 0 null null 0);
    }

    &-container {
        display: flex;
        flex-direction: column;
        margin-left: 27px;
    }

    input:checked ~ .radio-btn {
        &:after {
            @include position(absolute, 5px null null 5px);
            content: "";
            width: 10px;
            height: 10px;
            border-radius: 50%;
            background-color: #2c3e50;
        }
    }
}

.action-container {
    display: flex;
    flex-wrap: wrap;
}
.close-form-button {
    cursor: pointer;
    @include position(absolute, 14px 14px null null);
}
.btn {
    width: 280px;
    height: 45px;
    font-size: 18px;
    color: white;
    background-color: #ce596b;
    border: none;
    border-radius: 5px;
    font-family: $font-family-bold;
    letter-spacing: 0.1px;
    margin-bottom: 20px;

    &:nth-of-type(1) {
        margin-right: 18px;
    }

    &-delete {
        background: #fff;
        width: 215px;
        height: 45px;
        border: 1px solid #ce596b;
        color: #ce596b;
    }
}

.tags {
    &-container {
        position: relative;
        padding-top: 1em;

        > span {
            font-size: 100%;
            color: #6a6d71;
            transition: all 0.2s;
            cursor: text;
            z-index: $z-index-top;
            @include position(absolute, 11px null null 0);

            &.focused {
                font-size: 75%;
                top: -4px;
            }
        }
    }
}

::v-deep .vue-tags-input {
    max-width: unset;

    .ti-input {
        border: 0;
        border-bottom: 1px solid #d8dce0;
        padding: 0 0 3px 0;
    }

    .ti-new-tag-input {
        font-size: 16px;

        &-wrapper {
            padding: 0;
            margin: 0;
        }
    }

    .ti-tag {
        color: #2c3e50;
        font-size: 13px;
        height: 20px;
        border-radius: 10px;
        background-color: #f4f5f7;
    }

    .ti-icon-close {
        width: 10px;
        height: 10px;
        background: url("../../assets/icon-close.svg") no-repeat center center;
        background-size: 10px;
        background-color: transparent;

        &:before {
            content: "";
        }
    }

    .ti-deletion-mark {
        background: #ce596b !important;
        color: white;

        .ti-icon-close {
            background: url("../../assets/icon-close-white.svg") no-repeat center center;
            background-size: 10px;
        }
    }
}
::v-deep .has-float-label input {
    padding-bottom: 4px;
}
</style>
