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
                    <div>
                        <div class="name" v-if="projectUser.user.name != projectUser.user.email">{{ projectUser.user.name }}</div>
                        <div class="email">{{ projectUser.user.email }}</div>
                    </div>
                </div>

                <div v-if="projectUser.user.id !== user.id && !projectUser.invited" class="cell role">
                    <SelectField
                        :options="roleOptions"
                        :selected-label="projectUser.role"
                        @input="onEditMemberRole($event, projectUser.user.email)"
                    />
                </div>
                <div v-else class="cell">
                    <div>{{ projectUser.role }}</div>
                </div>

                <div class="cell" v-if="edited.email === projectUser.user.email">
                    <button class="invite-button" v-on:click="submitEditMemberRole()">Edit role</button>
                </div>
                <div class="cell invite-status" v-else>
                    <template v-if="projectUser.invited">
                        Invite pending
                    </template>
                </div>

                <div class="cell">
                    <img
                        v-if="projectUser.user.id !== user.id"
                        alt="Remove user"
                        src="@/assets/icon-close-bold.svg"
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
                    <SelectField :options="roleOptions" v-model="form.selectedRole" :selected-label="'Select Role'" />
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
                    disabled: true,
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
            edited: {
                email: null,
                role: null,
            },
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
        onEditMemberRole(role, email) {
            this.edited = {
                role,
                email,
            };
        },
        submitEditMemberRole() {
            const payload = {
                projectId: this.displayProject.id,
                role: this.edited.role,
                email: this.edited.email,
            };
            return this.$store.dispatch(ActionTypes.PROJECT_EDIT_ROLE, payload).then((response) => {
                this.edited.email = null;
                this.edited.role = null;
            });
        },
        sendInvite(this: any) {
            if (this.checkEmail()) {
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

<style scoped lang="scss">
@import "../../scss/mixins";

.manage-team-container {
    margin-top: 25px;
    display: flex;
    flex-direction: column;
    border: 1px solid #d8dce0;
    border-radius: 1px;
    background: white;
    padding: 25px 20px;

    @include bp-down($xs) {
        padding: 15px 10px 40px;
    }
}

.manage-team-container .section-heading {
    font-size: 20px;
    font-weight: 500;
    margin-bottom: 20px;

    @include bp-down($xs) {
        margin-bottom: 5px;
    }
}
.invite-button {
    width: 80px;
    height: 28px;
    border-radius: 3px;
    border: 1px solid #cccdcf;
    font-size: 14px;
    font-family: $font-family-bold;
    background-color: #ffffff;
    cursor: pointer;
}
.user-row {
    display: grid;
    font-size: 13px;
    grid-template-columns: 3fr 2fr 2fr 0.5fr;
    border-bottom: 1px solid rgb(215, 220, 225);
    padding: 10px 0;
    align-items: center;
    position: relative;

    &:last-of-type {
        margin-bottom: 20px;

        .cell {
            &:nth-of-type(1),
            &:nth-of-type(2) {
                @include bp-down($sm) {
                    padding-right: 85px;
                }
            }

            &:nth-of-type(3) {
                @include bp-down($sm) {
                    @include position(absolute, 50% 0 null null);
                    transform: translateY(-50%);
                }
            }
        }
    }

    @include bp-down($lg) {
        grid-template-columns: 2fr 1fr 1fr 0.5fr;
    }

    @include bp-down($sm) {
        grid-template-columns: repeat(auto-fit, 100%);
        padding: 3px 0 8px;
    }

    .cell {
        flex-wrap: wrap;

        &:nth-of-type(2) {
            color: #6a6d71;
        }

        &:nth-of-type(4) {
            justify-content: flex-end;
            padding: 0 20px;

            @include bp-down($sm) {
                padding: 0;
                @include position(absolute, 14px 0 null null);
            }
        }
    }
}

.cell-heading {
    font-size: 14px;
    font-family: $font-family-bold;

    &:nth-of-type(2) {
        @include bp-down($sm) {
            display: none;
        }
    }
}
.user-row .cell {
    @include flex(center);
    line-height: 1.23;

    &:nth-of-type(2),
    &:nth-of-type(3) {
        @include bp-down($sm) {
            padding-left: 41px;
        }
    }
}
.invite-status {
    color: #0a67aa;
    font-family: $font-family-normal !important;
}
.cell .text-input {
    border: none;
    border-radius: 5px;
    color: #818181;
    padding: 4px 0 4px 42px;
    font-size: 13px;
    display: flex;
    flex: 1;
    padding-right: 10px;
}
.cell .validation-error {
    color: #c42c44;
    display: block;

    @include bp-down($xs) {
        padding: 0 0 8px 42px;
    }
}
.cell .remove-button {
    cursor: pointer;
}
::v-deep .cell.role {
    margin-right: 1em;

    select {
        max-width: 150px;
        font-size: 13px;
    }
}
.users-container .user-icon {
    float: left;
}
.name {
    font-size: 14px;
}
.email {
    color: #818181;
}
</style>
