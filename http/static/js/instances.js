import {InstanceApi} from "./api/instance.js";
import {CheckBoxTop} from "./com/utils.js";


export class Instances {
    // {
    //   id: '#instances'
    // }
    constructor(props) {
        this.id = props.id;
        this.instanceOn = new InstanceOn(props);
        this.instances = this.instanceOn.uuids;

        // register buttons's click.
        $(`${this.id} #console`).on("click", this.instances, function (e) {
            let props = {uuids: e.data.store, passwd: {}};
            e.data.store.forEach(function (v) {
                props.passwd[v] = $(`input[data=${v}]`).attr('passwd');
            });
            new InstanceApi(props).console();
        });
        $(`${this.id} #start, ${this.id} #more-start`).on("click", this.instances, function (e) {
            new InstanceApi({uuids: e.data.store}).start();
        });
        $(`${this.id} #more-shutdown`).on("click", this.instances, function (e) {
            new InstanceApi({uuids: e.data.store}).shutdown();
        });
        $(`${this.id} #more-reset`).on("click", this.instances, function (e) {
            new InstanceApi({uuids: e.data.store}).reset();
        });
        $(`${this.id} #more-suspend`).on("click", this.instances, function (e) {
            new InstanceApi({uuids: e.data.store}).suspend();
        });
        $(`${this.id} #more-resume`).on("click", this.instances, function (e) {
            new InstanceApi({uuids: e.data.store}).resume();
        });
        $(`${this.id} #more-destroy`).on("click", this.instances, function (e) {
            new InstanceApi({uuids: e.data.store}).destroy();
        });
        $(`${this.id} #more-remove`).on("click", this.instances, function (e) {
            new InstanceApi({uuids: e.data.store}).remove();
        });
    }

    create(data) {
        new InstanceApi().create(data);
    }
}


export class InstanceOn {
    // {
    //   id: '#instances'
    // }
    constructor(props) {
        this.id = props.id;
        this.uuids = {store: [], id: this.id};

        let change = this.change;
        let record = this.uuids;

        new CheckBoxTop({
            one: `${this.id} #on-one`,
            all: `${this.id} #on-all`,
            change: function(e) {
                change(record, e);
            }
        });

        // disabled firstly.
        change(record, this.uuids);
    }

    change(record, from) {
        record.store = from.store;
        console.log(record.store);

        if (from.store.length == 0) {
            $(`${record.id} #start`).addClass('disabled');
            $(`${record.id} #console`).addClass('disabled');
            $(`${record.id} #shutdown`).addClass('disabled');
            $(`${record.id} #more`).addClass('disabled');
        } else {
            $(`${record.id} #start`).removeClass('disabled');
            $(`${record.id} #console`).removeClass('disabled');
            $(`${record.id} #shutdown`).removeClass('disabled');
            $(`${record.id} #more`).removeClass('disabled');
        }
        if (from.store.length != 1) {
            $(`${record.id} #edit`).addClass('disabled');
        } else {
            $(`${record.id} #edit`).removeClass('disabled');
        }
    }
}