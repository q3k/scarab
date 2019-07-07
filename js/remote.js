import common_pb from 'proto/common/common_js_proto/proto/common/common_pb.js';
import common_grpc from 'proto/common/common_js_proto/proto/common/common_grpc_web_pb.js';

export const Remote = function() {
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

Remote.prototype.state = async function() {
    const req = new common_pb.DefinitionsRequest();
    return await this._unary('state', req);
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
