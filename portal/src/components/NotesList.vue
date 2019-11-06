<template>
    <div id="notes-list-container" v-if="this.station">
        <div id="notes-title">Notes & Comments</div>
        <div v-for="note in notes" v-bind:key="note.id" class="note-container">
            <div class="top-line">
                <div class="creator">{{ note.creator }}</div>
                <div class="date">on {{ getDate(note.created) }}</div>
            </div>
            <div class="content-line">
                {{ note.note }}
            </div>
        </div>
        <div v-if="this.notes.length == 0">
            <p>Check back soon for notes and comments!</p>
        </div>
    </div>
</template>

<script>
import FKApi from "../api/api";

export default {
    name: "NotesList",
    props: ["station", "selectedSensor"],
    watch: {
        station() {
            if (this.station) {
                this.getFieldNotes();
            }
        }
    },
    data: () => {
        return {
            notes: []
        };
    },
    methods: {
        getFieldNotes() {
            const api = new FKApi();
            api.getFieldNotes(this.station).then(result => {
                this.notes = result.notes;
            });
        },
        getDate(strDate) {
            let d = new Date(strDate);
            return d.toLocaleDateString("en-US");
        }
    }
};
</script>

<style scoped>
#notes-list-container {
    float: left;
    width: 625px;
    margin: 10px;
    padding: 10px;
    border: 1px solid rgb(230, 230, 230);
}
#notes-title {
    margin: 10px 0;
    padding-bottom: 10px;
    border-bottom: 1px solid rgb(200, 200, 200);
}
.note-container {
    margin: 15px 0;
    padding-bottom: 10px;
    border-bottom: 1px solid rgb(230, 230, 230);
}
.creator {
    font-weight: bold;
    display: inline-block;
}
.date {
    font-size: 14px;
    margin-left: 6px;
    display: inline-block;
}
</style>
