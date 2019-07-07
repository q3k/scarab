import Vue from 'vue/dist/vue.esm.browser.js';
import VueRouter from 'vue-router/dist/vue-router.esm.browser.js';

import common_pb from 'proto/common/common_js_proto/proto/common/common_pb.js';

import { Remote } from './remote.js';
import { ViewIndex, ViewJobDefinition } from './components.js';

Vue.use(VueRouter);

const ModalState = Object.freeze({
        IDLE: Symbol("IDLE"),
        CREATE_JOB_SELECT_TYPE: Symbol("CREATE_JOB_SELECT_TYPE"),
        CREATE_JOB_INPUT_PARAMETERS: Symbol("CREATE_JOB_INPUT_PARAMETERS"),
        CREATE_JOB_START: Symbol("CREATE_JOB_START"),
        LOG: Symbol("LOG"),
});

const store = {
    state: {
        modalState: ModalState.IDLE,
        log: "",

        jobTypes: {},
        creatingJobName: "",
        creatingJobFieldErrors: {},

        statistics: [],
    },

    remote: new Remote(),

    logEntry(line) {
        this.state.log += line;
        this.state.log += "\n";
    },

    async tick() {
        const res = await this.remote.state();
        const statistics = res.getPerDefinitionStatisticsList();
        this.state.statistics = statistics.map((s) => {
            return {
                name: s.getDefinitionName(),
                description: s.getDefinitionDescription(),
                present: s.getJobsPresent(),
            };
        }).sort((a, b) => {
            if (a.name < b.name) {
                return -1;
            }
            if (a.name > b.name) {
                return 1;
            }
            return 0;
        });
    },

    idle() {
        this.state.modalState = ModalState.IDLE;
    },

    async jobSelect() {
        const definitions = await this.remote.definitions();
        const jobs = definitions.getJobsList();
        for (const job of jobs) {
            this.state.jobTypes[job.getName()] = job;
        }

        this.state.creatingJobName = "";
        this.state.creatingJobFieldErrors = {};
        this.state.modalState = ModalState.CREATE_JOB_SELECT_TYPE;
    },

    async jobInputParameters(jobName) {
        if (jobName === "") {
            return;
        }
        this.state.creatingJobName = jobName;
        this.state.modalState = ModalState.CREATE_JOB_INPUT_PARAMETERS;
    },

    async jobStart(fieldValues) {
        const job = this.state.jobTypes[this.state.creatingJobName];

        let errors = new Map();
        let fields = new Map();

        // Check fields.
        for (const argument of job.getArgumentsList()) {
            const name = argument.getName();
            let value = fieldValues[name] || "";
            let hasError = false;
            for (const validator of argument.getValidatorList()) {
                if (validator === common_pb.ArgumentDefinition.Validator.VALIDATOR_MUST_BE_SET) {
                    if (value === "") {
                        if (argument.getType() !== common_pb.ArgumentDefinition.Type.TYPE_BOOL) {
                            errors.set(name, "must be set");
                            hasError = true;
                            continue;
                        } else {
                            value = "false";
                        }
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
            for (const [k, v] of fields) {
                console.log(k, v);
                this.logEntry(`    ${k}: "${v}"`);
            }
        }
        this.logEntry("");

        // Create job on Scarab.

        this.state.modalState = ModalState.LOG;

        let res = undefined;
        try {
            res = await this.remote.create(this.state.creatingJobName, fields);
            let job_id = res.getJobId();
            this.logEntry(`Created new job, ID ${job_id}`);
        } catch (err) {
            this.logEntry(`Could not create job: ${err}`)
            return;
        }
    },
};

const router = new VueRouter({
    routes: [
        { path: "/", component: ViewIndex },
        { path: "/job/definition/:name", component: ViewJobDefinition },
    ],
})

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
            return this.state.modalState === ModalState.CREATE_JOB_SELECT_TYPE;
        },
        showCreateJobInputParameters: function() {
            return this.state.modalState === ModalState.CREATE_JOB_INPUT_PARAMETERS;
        },
        showLog: function() {
            return this.state.modalState === ModalState.LOG;
        },
    },
    mounted: function() {
        store.tick();
        setInterval(() => {
            store.tick();
        }, 5000);
    },
    router: router,
});
