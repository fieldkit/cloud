<template>
    <div class="row-section manage-team-container">
        <vue-confirm-dialog />

        <div class="section-heading">Manage Team</div>
        <div class="users-container">
            <div class="user-row">
                <div class="cell-heading">Members ({{ displayProject.users.length }})</div>
                <div class="cell-heading">Role</div>
                <div class="cell-heading"></div>
                <div class="cell"></div>
            </div>
            <div class="user-row" v-for="projectUser in displayProject.users" v-bind:key="projectUser.user.email">
                <div class="cell">
                    <UserPhoto :user="projectUser.user" />
                    <div class="invite-name">
                        <div v-if="projectUser.user.name != projectUser.user.email">{{ projectUser.user.name }}</div>
                        <div class="email">{{ projectUser.user.email }}</div>
                    </div>
                </div>
                <div class="cell">{{ projectUser.role }}</div>
                <div class="cell invite-status">Invite {{ projectUser.membership.toLowerCase() }}</div>
                <div class="cell">
                    <img
                        alt="Remove user"
                        src="@/assets/close-icon.png"
                        class="remove-button"
                        :data-user="projectUser.user.id"
                        v-on:click="removeUser(projectUser)"
                    />
                </div>
            </div>
            <div class="user-row">
                <div class="cell">
                    <input
                        class="text-input"
                        placeholder="New member email"
                        keyboardType="email"
                        autocorrect="false"
                        autocapitalizationType="none"
                        v-model="form.inviteEmail"
                    />
                    <div class="validation-errors" v-if="$v.form.inviteEmail.$error || form.inviteDuplicate">
                        <span class="validation-error" v-if="!$v.form.inviteEmail.required">
                            Email is a required field.
                        </span>
                        <span class="validation-error" v-if="!$v.form.inviteEmail.email">
                            Must be a valid email address.
                        </span>
                        <span class="validation-error" v-if="form.inviteDuplicate">
                            This user is already invited.
                        </span>
                    </div>
                </div>
                <div class="cell role">
                    <SelectField :options="roleOptions" v-model="form.selectedRole" />
                </div>
                <div class="cell">
                    <button class="invite-button" v-on:click="sendInvite">Invite</button>
                </div>
                <div class="cell"></div>
            </div>
        </div>
    </div>
</template>

<script lang="ts">
import Vue from "vue";
import UserPhoto from "@/views/shared/UserPhoto.vue";
import SelectField from "@/views/shared/SelectField.vue";
import { required, email } from "vuelidate/lib/validators";

import * as ActionTypes from "@/store/actions";
import FKApi from "@/api/api";

export default Vue.extend({
    name: "TeamManager",
    components: {
        UserPhoto,
        SelectField,
    },
    props: {
        displayProject: {
            type: Object,
            required: true,
        },
    },
    data: () => {
        return {
            form: {
                inviteEmail: "",
                inviteDuplicate: false,
                selectedRole: -1,
            },
            roleOptions: [
                {
                    value: -1,
                    label: "Select Role",
                },
                {
                    value: 0,
                    label: "Member",
                },
                {
                    value: 1,
                    label: "Administrator",
                },
            ],
        };
    },
    validations: {
        form: {
            inviteEmail: {
                required,
                email,
            },
        },
    },
    methods: {
        checkEmail(this: any) {
            this.$v.form.$touch();
            return !(this.$v.form.$pending || this.$v.form.$error);
        },
        sendInvite(this: any) {
            if (this.checkEmail()) {
                console.log(this.form.selectedRole);
                if (this.form.selectedRole === -1) {
                    this.form.selectedRole = 0;
                }
                const payload = {
                    projectId: this.displayProject.id,
                    email: this.form.inviteEmail,
                    role: this.form.selectedRole,
                };
                return this.$store
                    .dispatch(ActionTypes.PROJECT_INVITE, payload)
                    .then(() => {
                        this.$v.form.$reset();
                        this.form.inviteEmail = "";
                        this.form.inviteDuplicate = false;
                    })
                    .catch(() => {
                        // TODO: Move this to vuelidate.
                        this.form.inviteDuplicate = true;
                    });
            }
        },
        removeUser(this: any, projectUser) {
            return this.$confirm({
                message: `Are you sure you want to remove this team member?`,
                button: {
                    no: "No",
                    yes: "Yes",
                },
                callback: (confirm) => {
                    if (confirm) {
                        const payload = {
                            projectId: this.displayProject.id,
                            email: projectUser.user.email,
                        };
                        return this.$store.dispatch(ActionTypes.PROJECT_REMOVE, payload).then(() => {
                            this.$v.form.$reset();
                            this.form.inviteEmail = "";
                        });
                    }
                },
            });
        },
    },
});
</script>

<style scoped>
.manage-team-container {
    margin-top: 20px;
    display: flex;
    flex-direction: column;
    border: 2px solid #d8dce0;
    background: white;
    padding: 20px;
}

.manage-team-container .section-heading {
    font-weight: bold;
    margin-top: 10px;
    margin-bottom: 20px;
}
.invite-button {
    width: 80px;
    height: 28px;
    border-radius: 3px;
    border: 1px solid #cccdcf;
    font-size: 14px;
    font-weight: bold;
    background-color: #ffffff;
    cursor: pointer;
}
.user-row {
    display: grid;
    font-size: 13px;
    grid-template-columns: 4fr 1fr 1fr 1fr;
    border-bottom: 1px solid rgb(215, 220, 225);
    padding: 10px 0;
    align-items: center;
}
.cell .invite-status {
    color: #0a67aa;
    font-weight: 600;
}
.cell-heading {
    font-size: 14px;
    font-weight: bold;
}
.user-row .cell {
    text-align: left;
}
.cell .text-input {
    border: none;
    border-radius: 5px;
    font-size: 15px;
    padding: 4px 0 4px 8px;
}
.cell .validation-error {
    color: #c42c44;
    display: block;
}
.cell .remove-button {
    margin: 12px 0 0 0;
    float: right;
    cursor: pointer;
}
.cell .invite-name {
    display: inline-block;
    vertical-align: top;
    margin-top: 20px;
}
.cell.role {
    margin-right: 1em;
}
.users-container .user-icon {
    float: left;
}
</style>
