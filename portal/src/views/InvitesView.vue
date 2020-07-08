<template>
    <StandardLayout>
        <div id="loading" v-if="loading">
            <img alt="" src="../assets/progress.gif" />
        </div>
        <div class="main-panel" v-show="!loading">
            <div class="view-user">
                <div id="user-name">Hi, {{ this.user.name }}</div>
                <div class="notification" v-if="invalidToken">
                    Sorry, that invite link appears to have been already used or is invalid.
                </div>
                <div v-if="pending.length == 0" class="invite-heading">You have no pending invites.</div>
                <div v-if="pending.length > 0" class="invite-heading">You've been invited to the following projects:</div>
                <div v-for="invite in pending" v-bind:key="invite.id" class="project-row">
                    <div class="project-name">{{ invite.project.name }}</div>
                    <div class="accept-link" :data-id="invite.id" v-on:click="accept">Accept</div>
                    <div class="decline-link" :data-id="invite.id" v-on:click="decline">Decline</div>
                </div>
                <div v-for="invite in resolved" v-bind:key="invite.id" class="project-row">
                    <div class="project-name">{{ invite.project.name }}</div>
                    <div class="status">{{ invite.status }}</div>
                </div>
            </div>
        </div>
    </StandardLayout>
</template>

<script>
import StandardLayout from "./StandardLayout";
import FKApi from "../api/api";
import Config from "../secrets";

export default {
    name: "InvitesView",
    components: {
        StandardLayout,
    },
    data: () => {
        return {
            baseUrl: Config.API_HOST,
            user: { name: "" },
            projects: [],
            stations: [],
            pending: [],
            resolved: [],
            isAuthenticated: false,
            noCurrentUser: false,
            loading: false,
            invalidToken: false,
        };
    },
    async beforeCreate() {
        this.api = new FKApi();

        if (this.$route.query.token) {
            this.api.getInvitesByToken(this.$route.query.token).then(
                (result) => {
                    this.pending = result.pending;
                },
                () => {
                    this.invalidToken = true;
                }
            );
        }

        this.api
            .getCurrentUser()
            .then((user) => {
                this.user = user;
                this.isAuthenticated = true;
                this.api.getUserProjects().then((projects) => {
                    if (projects && projects.projects.length > 0) {
                        this.projects = projects.projects;
                    }
                });
                this.api.getStations().then((s) => {
                    this.stations = s.stations;
                });
            })
            .catch(() => {
                this.loading = false;
                this.noCurrentUser = true;
            });
    },
    methods: {
        goBack() {
            window.history.length > 1 ? this.$router.go(-1) : this.$router.push("/");
        },
        showStation(station) {
            this.$router.push({ name: "viewStation", params: { id: station.id } });
        },
        accept(event) {
            const inviteId = event.target.getAttribute("data-id");
            this.api.acceptInvite(inviteId).then(() => {
                const index = this.pending.findIndex((p) => {
                    return p.id == inviteId;
                });
                if (index > -1) {
                    const invite = this.pending.splice(index, 1)[0];
                    invite.status = "Accepted";
                    this.resolved.push(invite);
                }
            });
        },
        decline(event) {
            const inviteId = event.target.getAttribute("data-id");
            this.api.declineInvite(inviteId).then(() => {
                const index = this.pending.findIndex((p) => {
                    return p.id == inviteId;
                });
                if (index > -1) {
                    const invite = this.pending.splice(index, 1)[0];
                    invite.status = "Declined";
                    this.resolved.push(invite);
                }
            });
        },
    },
};
</script>

<style scoped>
#account-heading {
    font-weight: bold;
    font-size: 24px;
    float: left;
    margin: 15px 0 0 15px;
}
#loading {
    width: 100%;
    height: 100%;
    background-color: rgba(255, 255, 255, 0.65);
    text-align: center;
}
.no-user-message {
    float: left;
    font-size: 20px;
    margin: 40px 0 0 40px;
}
.show-link {
    text-decoration: underline;
}
#user-name {
    font-size: 24px;
    font-weight: bold;
    margin: 30px 15px 0 20px;
}
.invite-heading {
    float: left;
    clear: both;
    font-size: 18px;
    margin: 30px 0 0 20px;
}
.project-row {
    float: left;
    clear: both;
    min-width: 600px;
    padding: 8px;
    margin: 30px 0 0 20px;
    border-bottom: 2px solid #d8dce0;
}
.project-name {
    font-size: 22px;
    font-weight: bold;
    margin: 0 20px 0 0;
    display: inline-block;
}
.submit-btn {
    width: 100px;
    height: 30px;
    color: white;
    font-size: 18px;
    font-weight: 600;
    background-color: #ce596b;
    border: none;
    border-radius: 5px;
    margin: 0 20px 0 20px;
    cursor: pointer;
}
.accept-link,
.decline-link {
    margin: 0 20px;
    display: inline-block;
    font-size: 18px;
    font-weight: 600;
    text-decoration: underline;
    cursor: pointer;
}
.accept-link {
    color: #0a67aa;
}
.decline-link {
    color: #6b6d6f;
}
.status {
    margin: 0 20px;
    display: inline-block;
    font-size: 18px;
    font-weight: 600;
}
.notification {
    margin: 20px;
    padding: 20px;
    background-color: #f8d7da;
    border: 2px;
    border-radius: 4px;
}
</style>
