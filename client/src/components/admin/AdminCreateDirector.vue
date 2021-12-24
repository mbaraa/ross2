<template>
    <v-dialog
        max-height="400"
        max-width="-40"
        transition="dialog-bottom-transition"
        v-model="dialog"
        scrollable>
        <template v-slot:activator="{on, attrs}">
            <div v-bind="attrs" v-on="on" class="main3 bg-grey-darken-3" @click="dialog = true" title="just click it!">
                <h1>Add Director!</h1>
                <v-divider/>
                <FontAwesomeIcon style="font-size: 3em" :icon="{prefix:'fas', iconName:'plus'}"/>
            </div>
        </template>

        <v-card elevation="16" class="contestForm">
            <v-card-title>
                <span class="text-h4">Add Director</span>
            </v-card-title>

            <v-text-field label="Email Address" v-model="director.user.email"/>

            <v-btn class="bg-red" @click="dialog = false">
                Close
            </v-btn>&nbsp;
            <v-btn class="bg-blue" @click="createDirector()">
                Create
            </v-btn>
        </v-card>
    </v-dialog>
</template>

<script lang="ts">
import {defineComponent} from "vue";
import {FontAwesomeIcon} from "@fortawesome/vue-fontawesome";
import {library} from "@fortawesome/fontawesome-svg-core";
import {faPlus} from "@fortawesome/free-solid-svg-icons";
import Organizer from "@/models/Organizer";
import AdminRequests from "@/utils/requests/AdminRequests";

library.add(faPlus);

export default defineComponent({
    name: "AdminCreateDirector",
    components: {
        FontAwesomeIcon
    },
    data() {
        return {
            dialog: false,
            director: new Organizer(),
        }
    },
    methods: {
        async createDirector() {
            const resp = await AdminRequests.createDirector(this.director);
            if (!resp.ok) {
                window.alert(await resp.text());
                return;
            }
            window.alert("director was created successfully!");
            window.location.reload();
        },
    }
});
</script>

<style scoped>
.main3 {
    color: white;
    text-align: center;
    width: 350px;
    margin: 0 auto;
    height: auto;
    border-radius: 5px;
    padding: 5px;
    cursor: pointer;
    font-family: 'Ropa Sans', sans-serif;
}

.contestForm {
    padding: 10px;
    margin: 0 auto;
    width: 400px;
    overflow-y: auto;
}
</style>

