<template>
    <v-dialog
        max-height="400"
        max-width="-40"
        transition="dialog-bottom-transition"
        v-model="dialog"
        scrollable=""
    >
        <template v-slot:activator="{ on, attrs }">
            <div
                v-bind="attrs"
                v-on="on"
                @click="openDialog()"
                style="display: inline"
            >
                <v-btn icon class="bg-cyan" title="generate teams info posts">
                    <FontAwesomeIcon class="text-white" :icon="{ prefix: 'fas', iconName: 'images' }"/>&nbsp;
                </v-btn>
            </div>
        </template>

        <v-card elevation="16" class="teamForm">
            <v-card-title>
                <span class="text-h4">Generate Teams Posts</span>
            </v-card-title>

            <v-btn v-if="useSampleImage" @click="useSampleImage = false" class="bg-amber ma-1 text-white">upload and use
                image
            </v-btn>

            <div class="list" v-if="!useSampleImage">
                <a class="v-btn bg-purple text-white image pa-2 ma-1" href="/team_post_template.png" target="_blank">template
                    image sample</a>
                <a class="v-btn bg-purple text-white image pa-2 ma-1" href="/team_post_sample.png" target="_blank">sample
                    image for 3 members</a>
                <v-btn @click="useSampleImage = true" class="bg-purple ma-1 text-white">use sample image for 3 members
                </v-btn>

                <v-file-input show-size label="Template Image" prepend-icon @change="selectFile"/>

                <p>Team order props:</p>
                <v-text-field
                    required
                    type="number"
                    label="Start position X"
                    v-model="teamOrderProps.startPosition.x"
                />
                <v-text-field
                    required
                    type="number"
                    label="Start position Y"
                    v-model="teamOrderProps.startPosition.y"
                />
                <v-text-field required type="number" label="Font size" v-model="teamOrderProps.fontSize"/>
                <v-text-field required type="number" label="Width" v-model="teamOrderProps.width"/>

                <p>Team name props:</p>
                <v-text-field
                    required
                    type="number"
                    label="Start position X"
                    v-model="teamNameProps.startPosition.x"
                />
                <v-text-field
                    required
                    type="number"
                    label="Start position Y"
                    v-model="teamNameProps.startPosition.y"
                />
                <v-text-field required type="number" label="Font size" v-model="teamNameProps.fontSize"/>
                <v-text-field required type="number" label="Width" v-model="teamNameProps.width"/>

                <div v-for="(memberProp, i) in membersProps" :key="i">
                    <p>Member #{{ i + 1 }} name props:</p>
                    <v-text-field
                        required
                        type="number"
                        label="Start position X"
                        v-model="memberProp.startPosition.x"
                    />
                    <v-text-field
                        required
                        type="number"
                        label="Start position Y"
                        v-model="memberProp.startPosition.y"
                    />
                    <v-text-field required type="number" label="Font size" v-model="memberProp.fontSize"/>
                    <v-text-field required type="number" label="Width" v-model="memberProp.width"/>
                </div>
            </div>

            <div>
                <v-btn class="bg-red" @click="dialog = false">Close</v-btn>&nbsp;
                <v-btn class="bg-blue" @click="generatePosts()">Generate</v-btn>
            </div>
        </v-card>
    </v-dialog>
</template>

<script lang="ts">
import {defineComponent} from "vue";
import {FontAwesomeIcon} from "@fortawesome/vue-fontawesome";
import {library} from "@fortawesome/fontawesome-svg-core";
import {faImages} from "@fortawesome/free-solid-svg-icons";
import Contest from "@/models/Contest";
import OrganizerRequests from "@/utils/requests/OrganizerRequests";

library.add(faImages);

class Point2 {
    x?: number;
    y?: number;
}

class FieldProps {
    startPosition: Point2;
    width?: number;
    fontSize?: number;

    constructor() {
        this.startPosition = new Point2();
    }
}

export default defineComponent({
    name: "DirectorGenerateTeamsPosts",
    components: {
        FontAwesomeIcon
    },
    data() {
        return {
            dialog: false,
            templateImageFile: undefined,
            teamNameProps: new FieldProps(),
            teamOrderProps: new FieldProps(),
            membersProps: new Array<FieldProps>(),
            imageB64: "",
            useSampleImage: false,
        }
    },
    props: {
        contest: Contest
    },
    async mounted() {
        const membersCount = (await Contest.getContestFromServer(this.contest.id)).participation_conditions?.max_team_members;
        if (membersCount != undefined) {
            for (let i = 1; i <= membersCount; i++) {
                this.membersProps.push(new FieldProps());
            }
        }
    },
    methods: {
        openDialog() {
            this.dialog = true;
        },
        checkImageFile(): boolean {
            if (!this.useSampleImage && (this.templateImageFile === undefined || this.templateImageFile.type.indexOf("png") == -1)) {
                window.alert("select file of image/png type!");
                this.templateImageFile = undefined;
                return false;
            }
            return true;
        },
        selectFile(file: any) {
            this.templateImageFile = file.target.files[0];
            this.checkImageFile();
        },
        setB64Image(b: string | ArrayBuffer | null) {
            console.log("b64", b);
            this.imageB64 = b;
        },

        async readFile(): Promise<string | ArrayBuffer | null> {
            let res: string | ArrayBuffer | null = "";
            // 🙉🙊🙈 if it works it ain't stupid
            const toBase64 = () => new Promise((resolve, reject) => {
                const reader = new FileReader();
                reader.readAsDataURL(this.templateImageFile);
                reader.onload = () => {
                    resolve(reader.result);
                    res = reader.result;
                    return res;
                };
                reader.onerror = error => reject(error);
            });
            await toBase64();

            return res;
        },
        async generatePosts() {
            if (!this.checkImageFile()) {
                return;
            }

            const templateImageB64 = this.useSampleImage ? "" : await this.readFile();
            const zipFile = await OrganizerRequests.generateTeamsPosts({
                "contest": this.contest,
                "teamNameProps": this.teamNameProps,
                "teamOrderProps": this.teamOrderProps,
                "membersNamesProps": this.membersProps,
                "baseImage": this.useSampleImage ? "" : templateImageB64.substring(templateImageB64.indexOf(",") + 1),
            });

            const f = document.createElement("a");
            f.href = `data:application/zip;base64,${zipFile}`
            f.download = `${this.contest.name}'s teams posts.zip`;
            f.click();
        }
    }
});
</script>

<style scoped>
.teamForm {
    padding: 10px;
    margin: 0 auto;
    width: 450px;
    overflow-y: auto;
}

.list {
    overflow: hidden;
    overflow-y: scroll;
    height: 60vh;
}

.image {
    padding: 5px;
}
</style>

