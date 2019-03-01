// @flow weak
import _ from 'lodash';

import React, { Component } from 'react';
import { connect } from 'react-redux';

import { Link } from 'react-router-dom';

import { loadDevices, loadDeviceFiles } from '../actions';

import { API_HOST } from '../secrets';

import '../../css/files.css';
import '../../css/bootstrap.min.css';

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

    render() {
        const { files, deviceId } = this.props;

        if (deviceId) {
            const deviceFiles = files.filesByDevice[deviceId];

            if (!_.isObject(deviceFiles)) {
                return (<div></div>);
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
                                <th>File</th>
                                <th>Type</th>
                                <th>Size</th>
                                <th>Uploaded</th>
                            </tr>
                        </thead>
                        <tbody>
                            {deviceFiles.all.map(file => {
                                return (
                                    <tr key={file.file_id}>
                                        <td>{file.file_id}</td>
                                        <td>{file.file_type_name}</td>
                                        <td>{file.size}</td>
                                        <td>{file.time}</td>
                                        <td><a href={API_HOST + file.urls.csv}>CSV</a></td>
                                        <td><a href={API_HOST + file.urls.json}>JSON</a></td>
                                        <td><a href={API_HOST + file.urls.raw}>Raw</a></td>
                                    </tr>
                                );
                            })}
                        </tbody>
                    </table>
                    </div>
                </div>
            );
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
                {files.devices.map(device => {
                    return (
                        <tr className="device" key={device.device_id}>
                            <td><Link to={'/files/' + device.device_id}>{device.device_id}</Link></td>
                            <td>{device.number_of_files}</td>
                            <td>{device.last_file_time}</td>
                        </tr>
                    );
                })}
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
