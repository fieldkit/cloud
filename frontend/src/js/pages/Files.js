// @flow weak
import 'whatwg-fetch';

import _ from 'lodash';
import moment from 'moment';

import React, { Component } from 'react';
import { connect } from 'react-redux';

import { Link } from 'react-router-dom';

import { loadDevices, loadDeviceFiles } from '../actions';

import { FkPromisedApi } from '../api/calls';

import '../../css/files.css';
import '../../css/bootstrap.min.css';

class ConcatenatedFiles extends Component {
    state = {
        details: null
    }

    /*
    constructor() {
        super();
        this.timerToken = null;
        this.timer = this.timer.bind(this);
    }

    componentWillUnmount() {
        if (this.timerToken) {
            clearInterval(this.timerToken);
        }
    }

    timer() {
        const { deviceId } = this.props;

        FkPromisedApi.queryDeviceFilesDetails(deviceId).then(newDetails => {
            const { details: oldDetails } = this.state;

            console.log(newDetails, oldDetails);

            this.timerToken = setTimeout(this.timer, 1000);
        });
    }
    */

    componentWillMount() {
        const { deviceId } = this.props;

        FkPromisedApi.queryDeviceFilesDetails(deviceId).then(details => {
            this.setState({
                details
            });
        });
    }

    generate(url) {
        fetch(url, {}).then((res) => {
            /*
            if (this.timerToken) {
                clearInterval(this.timerToken);
            }
            this.timerToken = setTimeout(this.timer, 1000);
            */
        });
    }

    render() {
        const { details } = this.state;

        if (!_.isObject(details)) {
            return (<div>Loading</div>);
        }

        return (
            <div>
                <h4>Concatenated Files</h4>

                <div className="">
                    <div className="row">
                        <div className="col-sm-6">
                            <table className="table">
                                <tbody>
                                  {this.renderData()}
                                  {this.renderLogs()}
                                </tbody>
                            </table>
                        </div>
                        <div className="col-sm-6">
                        </div>
                    </div>
                </div>
            </div>
        );
    }

    renderData() {
        const { details } = this.state;
        const { device } = this.props;

        if (!_.isObject(details.files.data)) {
            return (
                <tr>
                  <th>Data</th>
                  <td><button className="btn btn-sm btn-success" onClick={() => this.generate(device.urls.data.generate)}>Generate</button></td>
                  <td colSpan="5">None</td>
                </tr>
            );
        }

        const dataMoment = moment(details.files.data.time);
        const dataTime = dataMoment.format('MMM Do YYYY, h:mm:ss a');
        const dataAgo = dataMoment.fromNow();

        return (
            <tr>
              <th>Data</th>
              <td><button className="btn btn-sm btn-success" onClick={() => this.generate(device.urls.data.generate)}>Generate</button></td>
              <td>{dataTime}</td>
              <td>{dataAgo}</td>
              <td>{details.files.data.size} bytes</td>
              <td><a target="_blank" href={device.urls.data.csv + "?dl=0"}>CSV</a> (<a href={device.urls.data.csv}>Download</a>)</td>
              <td><a target="_blank" href={device.urls.data.fkpb + "?dl=0"}>FKPB</a> (<a href={device.urls.data.fkpb}>Download</a>)</td>
            </tr>
        );
    }

    renderLogs() {
        const { details } = this.state;
        const { device } = this.props;

        if (!_.isObject(details.files.logs)) {
            return (
                <tr>
                  <th>Logs</th>
                  <td><button className="btn btn-sm btn-success" onClick={() => this.generate(device.urls.logs.generate)}>Generate</button></td>
                  <td colSpan="5">None</td>
                </tr>
            );
        }

        const logsMoment = moment(details.files.logs.time);
        const logsTime = logsMoment.format('MMM Do YYYY, h:mm:ss a');
        const logsAgo = logsMoment.fromNow();

        return (
            <tr>
              <th>Logs</th>
              <td><button className="btn btn-sm btn-success" onClick={() => this.generate(device.urls.logs.generate)}>Generate</button></td>
              <td>{logsTime}</td>
              <td>{logsAgo}</td>
              <td>{details.files.logs.size} bytes</td>
              <td><a target="_blank" href={device.urls.logs.csv + "?dl=0"}>CSV</a> (<a href={device.urls.logs.csv}>Download</a>)</td>
              <td><a target="_blank" href={device.urls.logs.fkpb + "?dl=0"}>FKPB</a> (<a href={device.urls.logs.fkpb}>Download</a>)</td>
            </tr>
        );
    }
}

class Files extends Component {
    props: {
        loadDevices: any => void,
    };

    componentWillMount() {
        const { loadDevices } = this.props;
        const { loadDeviceFiles, deviceId } = this.props;

        loadDevices();

        if (deviceId) {
            loadDeviceFiles(deviceId);
        }
    }

    componentDidUpdate(prevProps, prevState, snapshot) {
        const { loadDeviceFiles, deviceId } = this.props;

        if (deviceId && deviceId !== prevProps.deviceId) {
            loadDeviceFiles(deviceId);
        }
    }

    renderFiles(device, deviceFiles, deviceId) {
        return (
            <div className="files page container-fluid">
                <div className="header">
                    <div className="project-name"><Link to='/'>FieldKit Project</Link> / <Link to='/files'>Files</Link></div>
                </div>
                <div className="main-container">
                  <ConcatenatedFiles device={device} deviceId={deviceId} />

                  <h4> Individual Uploads </h4>

                  <div className="">
                    <table className="table table-striped table-sm">
                        <thead>
                            <tr>
                                <th>File</th>
                                <th>Type</th>
                                <th>Size</th>
                                <th>Uploaded</th>
                                <th>Links</th>
                            </tr>
                        </thead>
                        <tbody>
                            {deviceFiles.all.map(file => this.renderFile(file))}
                        </tbody>
                    </table>
                  </div>
                </div>
            </div>
        );
    }

    renderFile(file) {
        const time = moment(file.time).format('MMM Do YYYY, h:mm:ss a');

        return (
            <tr key={file.file_id}>
                <td>{file.file_id}</td>
                <td>{file.file_type_name}</td>
                <td>{file.size}</td>
                <td>{time}</td>
                <td>
                    <a target="_blank" href={file.urls.csv + "?dl=0"}>CSV</a> (<a href={file.urls.csv}>Download</a>) <span> | </span>
                    <a target="_blank" href={file.urls.json + "?dl=0"}>JSON</a> (<a href={file.urls.json}>Download</a>) <span> | </span>
                    <a href={file.urls.raw}>Raw</a>
                </td>
            </tr>
        );
    }

    renderDevice(device) {
        const time = moment(device.last_file_time).format('MMM Do YYYY, h:mm:ss a');;
        const places = _.join(_(device.locations.entries).map(le => le.places).uniq().filter(s => s !== null && s !== "").value(), ", ");

        return (
            <tr className="device" key={device.device_id}>
                <td><Link to={'/files/' + device.device_id}>{device.device_id}</Link></td>
                <td>{device.number_of_files}</td>
                <td>{time}</td>
                <td>{places}</td>
            </tr>
        );
    }

    render() {
        const { files, deviceId } = this.props;

        if (deviceId) {
            const deviceFiles = files.filesByDevice[deviceId];
            const device = _(files.devices).filter(d => d.device_id === deviceId).first();

            if (!_.isObject(device) || !_.isObject(deviceFiles)) {
                return (<div></div>);
            }

            return this.renderFiles(device, deviceFiles, deviceId);
        }

        return (
            <div className="files page container-fluid">
                <div className="header">
                    <div className="project-name"><Link to='/'>FieldKit Project</Link> / <Link to='/files'>Files</Link></div>
                </div>
                <div className="main-container">

                <table className="table table-striped table-sm">
                <thead>
                    <tr>
                        <th>ID</th>
                        <th>Files</th>
                        <th>Last File</th>
                        <th>Places</th>
                    </tr>
                </thead>
                <tbody>
                    {files.devices.map(device => this.renderDevice(device))}
                </tbody>
                </table>

                </div>
            </div>
        );
    }
}

const mapStateToProps = state => ({
    files: state.files
});

function showWhenReady(WrappedComponent, isReady) {
    return class extends React.Component {
        render() {
            if (!isReady(this.props)) {
                return <div>Loading</div>;
            }

            return <WrappedComponent {...this.props} />;
        }
    };
};

export default connect(mapStateToProps, {
    loadDevices,
    loadDeviceFiles
})(showWhenReady(Files, props => {
    return true;
}));
