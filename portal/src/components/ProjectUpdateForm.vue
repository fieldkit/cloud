<template>
    <div id="project-form-container">
        <div id="close-form-btn" v-on:click="closeForm">
            <img alt="Close" src="../assets/close.png" />
        </div>
        <h2>{{ this.formHeading }}</h2>
        <div class="outer-input-container">
            <div class="input-container">
                <span class="floating-label">Enter your project update below.</span>
                <textarea v-model="body" class="inputText" required="" />
            </div>
        </div>

        <div class="action-container">
            <button class="save-btn" v-if="formType == 'add'" v-on:click="addProjectUpdate">
                Add
            </button>
            <button class="save-btn" v-if="formType == 'update'" v-on:click="updateProjectUpdate">
                Update
            </button>
            <div v-if="formType == 'update'" class="delete-container" v-on:click="deleteProjectUpdate">
                <img alt="Delete" src="../assets/Delete.png" />
                Delete this update
            </div>
        </div>
    </div>
</template>

<script>
import FKApi from "../api/api";

export default {
    name: "ProjectUpdateForm",
    props: ["projectUpdate", "project"],
    data: () => {
        return {
            formType: "add",
            formHeading: "New Project Update",
            body: "",
        };
    },
    mounted() {
        if (!this.project) {
            return this.$router.push({ name: "projects" });
        }
        if (this.projectUpdate) {
            this.formHeading = "Edit Project Update for " + this.project.name;
            this.formType = "update";
            this.body = this.projectUpdate.body;
        } else {
            this.formType = "add";
            this.formHeading = "New Project Update for " + this.project.name;
            this.body = "";
        }
    },
    methods: {
        createParams() {
            const data = {
                projectId: this.project.id,
                body: this.body,
            };
            if (this.projectUpdate) {
                // handle edit instead of add...
                // data.updateId = this.projectUpdate.id
            }
            return data;
        },
        addProjectUpdate() {
            this.$emit("updating");
            const api = new FKApi();
            const data = this.createParams();
            api.addProjectUpdate(data).then(() => {
                this.$router.push({ name: "viewProject", params: { id: this.project.id } });
            });
        },
        updateProjectUpdate() {
            this.$emit("updating");
            const api = new FKApi();
            const data = this.createParams();
            api.updateProjectUpdate(data).then(() => {
                this.$router.push({ name: "viewProject", params: { id: this.project.id } });
            });
        },
        deleteProjectUpdate() {
            if (window.confirm("Are you sure you want to delete this project update?")) {
                const api = new FKApi();
                const params = {
                    projectId: this.project.id,
                    updateId: this.projectUpdate.id,
                };
                api.deleteProjectUpdate(params).then(() => {
                    this.$router.push({ name: "viewProject", params: { id: this.project.id } });
                });
            }
        },
        closeForm() {
            this.$router.push({ name: "viewProject", params: { id: this.project.id } });
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
.input-container {
    margin: auto;
    width: 100%;
    text-align: left;
}
textarea {
    min-height: 200px;
    padding: 6px;
    border: 1px solid rgb(215, 220, 225);
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
