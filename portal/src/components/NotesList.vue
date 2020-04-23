<template>
    <div id="notes-list-container">
        <!-- <div id="notes-title">Notes & Comments</div> -->
        <!-- temporarily commandeered for field notes -->
        <div id="notes-title">Field Notes</div>
        <div v-for="note in notes" v-bind:key="note.id" class="note-container">
            <div class="delete-link">
                <img alt="Info" src="../assets/Delete.png" :data-id="note.id" v-on:click="deleteNote" />
            </div>
            <div class="top-line" v-if="!note.title">
                <div class="creator">{{ note.creator }}</div>
                <div class="date">on {{ getDate(note.created) }}</div>
            </div>
            <div v-if="note.media_id">
                <img alt="Field Note image" :src="getImageUrl(note)" class="field-note-image" />
            </div>
            <div class="content-line">
                <span class="title">{{ note.title }}</span>
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
import { API_HOST } from "../secrets";

// TODO: api return names for categories rather than keys (?)
const categoryLabels = {
    locationName: "Location Name: ",
    studyObjective: "Study Objective: ",
    locationPurpose: "Location Purpose: ",
    siteCriteria: "Site Criteria: ",
    siteDescription: "Site Description: ",
};

export default {
    name: "NotesList",
    props: [],
    data: () => {
        return {
            baseUrl: API_HOST,
            notes: [],
            ids: [],
        };
    },
    methods: {
        updateNotes(ids) {
            this.ids = ids;
            this.notes = [];
            this.getFieldNotes();
        },
        getFieldNotes() {
            const api = new FKApi();
            this.ids.forEach(id => {
                api.getFieldNotes(id).then(result => {
                    this.notes = _.concat(
                        this.notes,
                        result.notes.map(n => {
                            n.stationId = id;
                            if (n.category_key != "default") {
                                n.title = categoryLabels[n.category_key];
                            }
                            return n;
                        })
                    );
                });
            });
        },
        getImageUrl(note) {
            return this.baseUrl + "/stations/" + note.stationId + "/field-note-media/" + note.media_id;
        },
        getDate(strDate) {
            let d = new Date(strDate);
            return d.toLocaleDateString("en-US");
        },
        deleteNote(event) {
            const fieldNoteId = event.target.getAttribute("data-id");
            if (window.confirm("Are you sure you want to delete this note?")) {
                const api = new FKApi();
                const note = this.notes.find(n => {
                    return n.id == fieldNoteId;
                });
                const params = {
                    stationId: note.stationId,
                    fieldNoteId: fieldNoteId,
                };
                api.deleteFieldNote(params).then(() => {
                    const index = this.notes.findIndex(n => {
                        return n.id == fieldNoteId;
                    });
                    if (index > -1) {
                        this.notes.splice(index, 1);
                    }
                });
            }
        },
    },
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
    font-weight: bold;
    margin: 10px 0;
    padding-bottom: 10px;
    border-bottom: 1px solid rgb(200, 200, 200);
}
.top-line {
    display: inline-block;
}
.field-note-image {
    max-width: 300px;
}
.note-container {
    margin: 15px 0;
    padding-bottom: 10px;
    border-bottom: 1px solid rgb(230, 230, 230);
}
.title {
    font-weight: bold;
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
.content-line {
    width: 600px;
    display: inline-block;
}
.delete-link {
    float: right;
    opacity: 0;
}
.delete-link:hover {
    opacity: 1;
}
</style>
