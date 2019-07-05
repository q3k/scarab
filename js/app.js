import Vue from 'vue/dist/vue.esm.browser.js';

import common_pb from 'proto/common/common_js_proto/proto/common/common_pb.js';
import common_grpc from 'proto/common/common_js_proto/proto/common/common_grpc_web_pb.js';

const Remote = function() {
    const url = `${window.location.protocol}//${window.location.hostname}:${window.location.port}`;
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
        CREATE_JOB_START: Symbol("CREATE_JOB_START"),
});

const store = {
    state: {
        state: State.IDLE,

        jobTypes: {},
        creatingJobName: "",
    },

    remote: new Remote(),

    idle: function() {
        this.state.state = State.IDLE;
    },
    jobSelect: async function() {
        const descriptions = await this.remote.descriptions();
        const jobs = descriptions.getJobsList();
        for (const job of jobs) {
            this.state.jobTypes[job.getName()] = job;
        }
        this.state.state = State.CREATE_JOB_SELECT_TYPE;
    },
    jobInputParameters: async function(jobName) {
        if (jobName === "") {
            return;
        }
        this.state.creatingJobName = jobName;
        this.state.state = State.CREATE_JOB_INPUT_PARAMETERS;
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
    data: () => { return {
        'selected': "",
    }; },
    props: {
        "jobTypes": {type: Object, default: {}},
    },
    template: `
    <div id="modal">
        <div id="modalContent">
            <h3>Select Job type...</h3>
            <select id="select" v-model="selected">
                <option disabled value="">Select...</option>
                <option v-for="(job, name) in jobTypes" :value="name">{{ job.getDescription() }}</option>
            </select>
            <vbutton v-on:click="$emit('ok', selected)" red :s="{ marginRight: 0, }">OK</vbutton>
            <vbutton v-on:click="$emit('close')">Cancel</vbutton>
        </div>
    </div>
    `,
});

Vue.component('modal-job-input-parameters', {
    data: () => { return {
        fields: {},
    }; },
    props: {
        "job": {type: common_pb.JobDefinition},
        "fieldErrors": {type: Object, default: () => { return {}; }},
    },
    methods: {
        getArguments: function() {
            let res = [];
            
            for (const argument of this.job.getArgumentsList()) {
                let arg = {
                    name: argument.getName(),
                    description: argument.getDescription(),
                    checkbox: argument.getType() === common_pb.ArgumentDefinition.Type.TYPE_BOOL,
                    mustBeSet: false,
                };
                for (const validator of argument.getValidatorList()) {
                    if (validator == common_pb.ArgumentDefinition.Validator.VALIDATOR_MUST_BE_SET) {
                        arg.mustBeSet = true;
                    }
                }
                res.push(arg);
            }

            return res;
        }
    },
    template: `
    <div id="modal">
        <div id="modalContent">
            <h3>{{ job.getDescription() }} ...</h3>
            <div class="fields">
                <template v-for="argument in getArguments()">
                    <label :for=argument.name>
                        {{ argument.description }}<span v-if="argument.mustBeSet" style="color: red;"> *</span>:
                    </label>
                    <input
                        v-if="argument.checkbox"
                        v-model=fields[argument.name]
                        type="checkbox"
                        :name=argument.name
                    />
                    <input
                        v-else
                        v-model=fields[argument.name]
                        :name=argument.name
                    />
                </template>
            </div>
            <vbutton v-on:click="$emit('ok', fields)" red :s="{ marginRight: 0, }">Create Job</vbutton>
            <vbutton v-on:click="$emit('close')">Cancel</vbutton>
        </div>
    </div>
    `,
});

const vm = new Vue({
    el: '#app',
    data: {
        state: store.state,
    },
    methods: {
        idle: () => store.idle(),
        jobSelect: () => store.jobSelect(),
        jobInputParameters: (jobName) => store.jobInputParameters(jobName),
        jobStart: (fields) => console.log(fields),
    },
    computed: {
        showCreateJobSelectType: function() {
            return this.state.state === State.CREATE_JOB_SELECT_TYPE;
        },
        showCreateJobInputParameters: function() {
            return this.state.state === State.CREATE_JOB_INPUT_PARAMETERS;
        },
    },
});
