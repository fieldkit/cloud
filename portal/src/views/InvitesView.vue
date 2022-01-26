<template>
    <StandardLayout>
        <div class="invites-view" v-if="user && !loading">
            <div id="user-name">Hi, {{ user.name }}</div>
            <div class="notification" v-if="invalidToken">Sorry, that invite link appears to have been already used or is invalid.</div>
            <div v-if="pending.length == 0" class="invite-heading">You have no pending invites.</div>
            <div v-if="pending.length > 0" class="invite-heading">You've been invited to the following projects:</div>
            <div v-for="invite in pending" v-bind:key="invite.id" class="project-row">
                <div class="project-name">{{ invite.project.name }}</div>
                <div class="accept-link" :data-id="invite.id" v-on:click="(ev) => accept(ev, invite)">Accept</div>
                <div class="decline-link" :data-id="invite.id" v-on:click="(ev) => decline(ev, invite)">Decline</div>
            </div>
        </div>
    </StandardLayout>
</template>

<script>
import StandardLayout from "./StandardLayout";
import * as ActionTypes from "@/store/actions";

export default {
    name: "InvitesView",
    components: {
        StandardLayout,
    },
    data: () => {
        return {
            pending: [],
            loading: false,
            invalidToken: false,
        };
    },
    computed: {
        user() {
            return this.$store.state.user.user;
        },
    },
    async beforeCreate() {
        if (this.$route.query.token) {
            this.$services.api.getInvitesByToken(this.$route.query.token).then(
                (result) => {
                    this.pending = result.pending;
                },
                () => {
                    this.invalidToken = true;
                    this.pending = [];
                }
            );
        }
    },
    methods: {
        goBack() {
            window.history.length > 1 ? this.$router.go(-1) : this.$router.push("/");
        },
        showStation(station) {
            this.$router.push({ name: "mapStation", params: { id: station.id } });
        },
        accept(ev, invite) {
            const token = this.$route.query.token;
            return this.$store.dispatch(ActionTypes.ACCEPT_PROJECT_INVITE, { id: invite.id, token: token }).then(
                () => {
                    this.pending = [];
                },
                () => {
                    this.invalidToken = true;
                }
            );
        },
        decline(ev, invite) {
            const token = this.$route.query.token;
            return this.$store.dispatch(ActionTypes.DECLINE_PROJECT_INVITE, { id: invite.id, token: token }).then(
                () => {
                    this.pending = [];
                },
                () => {
                    this.invalidToken = true;
                }
            );
        },
    },
};
</script>

<style scoped>
.invites-view {
    display: flex;
    flex-direction: column;
    height: 100%;
    background-color: #fcfcfc;
    padding: 40px;
    text-align: left;
}
.show-link {
    text-decoration: underline;
}
#user-name {
    font-size: 24px;
    font-weight: bold;
}
.invite-heading {
    float: left;
    clear: both;
    font-size: 18px;
}
.project-row {
    float: left;
    clear: both;
    min-width: 600px;
    margin: 30px 0 0 0;
    border-bottom: 2px solid var(--color-border);
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
    font-family: var(--font-family-bold);
    background-color: var(--color-secondary);
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
    padding: 20px;
    background-color: #f8d7da;
    border: 2px;
    border-radius: 4px;
}
</style>
