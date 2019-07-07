import Vue from 'vue/dist/vue.esm.browser.js';

import common_pb from 'proto/common/common_js_proto/proto/common/common_pb.js';
import common_grpc from 'proto/common/common_js_proto/proto/common/common_grpc_web_pb.js';

const Remote = function() {
    this.url = `${window.location.protocol}//${window.location.hostname}:${window.location.port}`;
    console.log("Remote URL: " + this.url);
    this.manage = new common_grpc.ManageClient(this.url);
}

Remote.prototype._unary = function(method, req) {
    return new Promise((resolve, reject) => {
        this.manage[method](req, {}, (err, response) => {
            if (err) {
                reject(`RPC Error: ${method}: ${err.message}`);
            } else {
                resolve(response);
            }
        });
    });
}

Remote.prototype.definitions = async function() {
    const req = new common_pb.DefinitionsRequest();
    return await this._unary('definitions', req);
}

Remote.prototype.create = async function(jobName, fields) {
    let req = new common_pb.CreateRequest();
    req.setJobDefinitionName(jobName);
    let args = [];
    for (const [k, v] of fields) {
        let arg = new common_pb.Argument();
        arg.setName(k);
        arg.setValue(v);
        args.push(arg);
    }
    req.setArgumentsList(args);

    return await this._unary('create', req);
}

const State = Object.freeze({
        IDLE: Symbol("IDLE"),
        CREATE_JOB_SELECT_TYPE: Symbol("CREATE_JOB_SELECT_TYPE"),
        CREATE_JOB_INPUT_PARAMETERS: Symbol("CREATE_JOB_INPUT_PARAMETERS"),
        CREATE_JOB_START: Symbol("CREATE_JOB_START"),
        LOG: Symbol("LOG"),
});

const store = {
    state: {
        state: State.IDLE,

        jobTypes: {},
        creatingJobName: "",
        creatingJobFieldErrors: {},
        log: "",
    },

    remote: new Remote(),

    logEntry(line) {
        this.state.log += line;
        this.state.log += "\n";
    },

    idle() {
        this.state.state = State.IDLE;
    },

    async jobSelect() {
        const definitions = await this.remote.definitions();
        const jobs = definitions.getJobsList();
        for (const job of jobs) {
            this.state.jobTypes[job.getName()] = job;
        }

        this.state.creatingJobName = "";
        this.state.creatingJobFieldErrors = {};
        this.state.state = State.CREATE_JOB_SELECT_TYPE;
    },

    async jobInputParameters(jobName) {
        if (jobName === "") {
            return;
        }
        this.state.creatingJobName = jobName;
        this.state.state = State.CREATE_JOB_INPUT_PARAMETERS;
    },

    async jobStart(fieldValues) {
        const job = this.state.jobTypes[this.state.creatingJobName];

        let errors = new Map();
        let fields = new Map();

        // Check fields.
        for (const argument of job.getArgumentsList()) {
            const name = argument.getName();
            const value = fieldValues[name] || "";
            let hasError = false;
            for (const validator of argument.getValidatorList()) {
                if (validator === common_pb.ArgumentDefinition.Validator.VALIDATOR_MUST_BE_SET) {
                    if (value === "" && argument.getType() !== common_pb.ArgumentDefinition.Type.TYPE_BOOL) {
                        errors.set(name, "must be set");
                        hasError = true;
                        continue;
                    }
                }
            }
            if (hasError) {
                continue;
            }
            fields.set(name, value);
        }

        if (errors.size > 0) {
            // A validation error occured.
            console.log("Validation errors: ", errors);
            this.state.creatingJobFieldErrors = Object.fromEntries(errors);
            return;
        }
        this.state.creatingJobFieldErrors = {};

        // Start logging...

        this.state.log = "";
        this.logEntry(`Creating job "${this.state.creatingJobName}" on ${this.remote.url}`);
        if (fields.size > 0) {
            this.logEntry(`With arguments:`);
            console.log("wut", fields);
            for (const [k, v] of fields) {
                console.log(k, v);
                this.logEntry(`    ${k}: "${v}"`);
            }
        }
        this.logEntry("");

        // Create job on Scarab.

        this.state.state = State.LOG;

        let res = undefined;
        try {
            res = await this.remote.create(this.state.creatingJobName, fields);
        } catch (err) {
            this.logEntry(`Could not create job: ${err}`)
            return;
        }
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
            const ad = common_pb.ArgumentDefinition;
            return this.job.getArgumentsList().map((argument) => {
                const validators = argument.getValidatorList();
                return {
                    name: argument.getName(),
                    description: argument.getDescription(),
                    checkbox: argument.getType() === ad.Type.TYPE_BOOL,
                    mustBeSet: validators.some((v) => v == ad.Validator.VALIDATOR_MUST_BE_SET),
                };
            });
        }
    },
    template: `
    <div id="modal">
        <div id="modalContent">
            <h3>{{ job.getDescription() }} ...</h3>
            <div class="fields">
                <template v-for="argument in getArguments()">
                    <label :for=argument.name>
                        {{ argument.description }}<span v-if="argument.mustBeSet && !argument.checkbox" style="color: red;"> *</span>:
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
                    <div
                        v-if="fieldErrors[argument.name] !== undefined"
                        class="error"
                    >{{ fieldErrors[argument.name] }}</div>
                </template>
            </div>
            <vbutton v-on:click="$emit('ok', fields)" red :s="{ marginRight: 0, }">Create Job</vbutton>
            <vbutton v-on:click="$emit('close')">Cancel</vbutton>
        </div>
    </div>
    `,
});

Vue.component('modal-log', {
    data: () => { return {
    }; },
    props: {
        "job": {type: common_pb.JobDefinition},
        "log": {type: String, default: "..."},
    },
    template: `
    <div id="modal">
        <div id="modalContent">
            <h3>{{ job.getDescription() }} ...</h3>
            <pre class="log">{{ log }}</pre>
            <vbutton v-on:click="$emit('close')">Close</vbutton>
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
        jobStart: (fields) => store.jobStart(fields),
    },
    computed: {
        showCreateJobSelectType: function() {
            return this.state.state === State.CREATE_JOB_SELECT_TYPE;
        },
        showCreateJobInputParameters: function() {
            return this.state.state === State.CREATE_JOB_INPUT_PARAMETERS;
        },
        showLog: function() {
            return this.state.state === State.LOG;
        },
    },
});
