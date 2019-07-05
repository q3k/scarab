import Vue from 'vue/dist/vue.esm.browser.js';

import common_pb from 'proto/common/common_js_proto/proto/common/common_pb.js';
import common_grpc from 'proto/common/common_js_proto/proto/common/common_grpc_web_pb.js';

let Remote = function() {
    let url = `${window.location.protocol}//${window.location.hostname}:${window.location.port}`;
    console.log("Remote URL: " + url);
    this.manage = new common_grpc.ManageClient(url);
}

Remote.prototype.descriptions = function() {
    let req = new common_pb.DefinitionsRequest();
    return new Promise((resolve, reject) => {
        this.manage.definitions(req, {}, (err, response) => {
            if (err) {
                reject(err.message);
            } else {
                resolve(response);
            }
        });
    })
}

const State = Object.freeze({
        IDLE: Symbol("IDLE"),
        CREATE_JOB_SELECT_TYPE: Symbol("CREATE_JOB_SELECT_TYPE"),
        CREATE_JOB_INPUT_PARAMETERS: Symbol("CREATE_JOB_INPUT_PARAMETERS"),
});

let store = {
    state: {
        state: State.IDLE,

        jobTypes: {},
    },

    remote: new Remote(),

    idle: function() {
        this.state.state = State.IDLE;
    },
    jobSelect: async function() {
        let descriptions = await this.remote.descriptions();
        let jobs = descriptions.getJobsList();
        for (let job of jobs) {
            this.state.jobTypes[job.getName()] = job.getDescription();
        }
        this.state.state = State.CREATE_JOB_SELECT_TYPE;
    },
};

Vue.component('vbutton', {
    data: () => {
        return { }
    },
    props: {
        red: { type: Boolean, default: false },
        s: { type: Object, default: function() { return {} } },
    },
    template: `
        <a href="#" v-bind:class="{ button: true, red: red}" v-bind:style=s v-on:click="$emit('click')"><slot></slot></a>
    `,
});

Vue.component('modal-job-create', {
    props: {
        "jobTypes": {type: Object, default: {}},
    },
    template: `
    <div id="modal">
        <div id="modalContent">
            <h3>Select Job type...</h3>
            <select id="select">
                <option v-for="(description, name) in jobTypes">{{ description }}</option>
            </select>
            <vbutton red :s="{ marginRight: 0, }">OK</vbutton>
            <vbutton v-on:click="$emit('close')">Cancel</vbutton>
        </div>
    </div>
    `,
});

let vm = new Vue({
    el: '#app',
    data: {
        state: store.state,
    },
    methods: {
        idle: () => store.idle(),
        jobSelect: () => store.jobSelect(),
    },
    computed: {
        showCreateJobSelectType: function() {
            return this.state.state == State.CREATE_JOB_SELECT_TYPE;
        }
    },
});
