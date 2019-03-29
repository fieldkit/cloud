// @flow weak
import 'whatwg-fetch';

import _ from 'lodash';
import moment from 'moment';

import React, { Component } from 'react';
import { connect } from 'react-redux';

import { Link } from 'react-router-dom';

import { loadDevices, loadDeviceFiles } from '../actions';

import { API_HOST } from '../secrets';

import '../../css/files.css';
import '../../css/bootstrap.min.css';

class ConcatenatedFiles extends Component {
    generate(url) {
        console.log(url);
        fetch(url, {}).then((res) => {
            console.log("DONE", res.body, res.headers);
        });
    }

    render() {
        const { device } = this.props;

        return (
            <div>
                <h4> Concatenated Files </h4>

                <div className="">
                    <div className="row">
                        <div className="col-sm-3">
                            <table className="table">
                            <tbody>
                              <tr>
                                <th>Data</th>
                                <td><button className="btn btn-sm btn-success" onClick={() => this.generate(device.urls.data.generate)}>Generate</button></td>
                                <td><a target="_blank" href={device.urls.data.csv}>CSV</a></td>
                                <td><a target="_blank" href={device.urls.data.csv}>FKPB</a></td>
                              </tr>
                              <tr>
                                <th>Logs</th>
                                <td><button className="btn btn-sm btn-success" onClick={() => this.generate(device.urls.logs.generate)}>Generate</button></td>
                                <td><a target="_blank" href={device.urls.logs.csv}>CSV</a></td>
                                <td><a target="_blank" href={device.urls.logs.csv}>FKPB</a></td>
                              </tr>
                            </tbody>
                            </table>
                        </div>
                        <div className="col-sm-9">
                        </div>
                    </div>
                </div>
            </div>
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
        const time = moment(file.time).format('MMM Do YYYY, h:mm:ss a');;

        return (
            <tr key={file.file_id}>
                <td>{file.file_id}</td>
                <td>{file.file_type_name}</td>
                <td>{file.size}</td>
                <td>{time}</td>
                <td>
                    <a target="_blank" href={API_HOST + file.urls.csv + "?dl=0"}>CSV</a> (<a href={API_HOST + file.urls.csv}>Download</a>) <span> | </span>
                    <a target="_blank" href={API_HOST + file.urls.json + "?dl=0"}>JSON</a> (<a href={API_HOST + file.urls.json}>Download</a>) <span> | </span>
                    <a href={API_HOST + file.urls.raw}>Raw</a>
                </td>
            </tr>
        );
    }

    renderDevice(device) {
        const time = moment(device.last_file_time).format('MMM Do YYYY, h:mm:ss a');;

        return (
            <tr className="device" key={device.device_id}>
                <td><Link to={'/files/' + device.device_id}>{device.device_id}</Link></td>
                <td>{device.number_of_files}</td>
                <td>{time}</td>
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
