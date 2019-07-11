import Vue from 'vue/dist/vue.esm.browser.js';

import common_pb from 'proto/common/common_js_proto/proto/common/common_pb.js';

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

Vue.component('job-definition-list', {
    props: {
        statistics: {type: Object, },
    },
    template: `
    <div class="row">
        <vbutton red v-on:click="$emit('create')" :s="{ marginBottom: '10px', }">Create Job</vbutton>
        <ul>
            <li v-for="(stat, name) in statistics">
                <router-link :to="'/job/definition/' + name" class="job" active-class="selected">{{ stat.description }}</router-link>
            </li>
        </ul>
    </div>
    `,
});

export const ViewIndex = {
    template: `
    <div class="row">
        <h3>Welcome to Scarab!</h3>
        <p>
            Select a Job type from the sidebar to show its progress here. Or, click 'Create Job' to start a new job.
        </p>
    </div>
    `
};

export const ViewJobDefinition = {
    props: {
        state: {type: Object, },
    },
    template: `
        <div>
            <div v-for="job in state.jobsPerDefinition[$route.params.name]" class="row">
                <h3>{{ job.getDefinition().getDescription() }} <small>(#{{ job.getId() }})</small></h3>
                <p>
                    Started with arguments:
                </p>
                <ul>
                    <li v-for="arg in job.getArgumentsList()"><b>{{ arg.getName() }}</b>: {{ arg.getValue() }}</li>
                </li>
            </div>
            <div v-if="(state.jobsPerDefinition[$route.params.name] || []).length === 0" class="row">
                <h3>No active jobs of type '{{ $route.params.name }}'.</h3>
                <p>
                    Select a job type from the left.
                </p>
            </div>
        </div>
    `,
    async beforeRouteEnter(to, from, next) {
        await store.fetchJobDefinitionData(to.params.name);
        next();
    },
};
