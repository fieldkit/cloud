// @flow weak
import 'whatwg-fetch';

import _ from 'lodash';
import moment from 'moment';
import prettyBytes from 'pretty-bytes';

import React, { Component } from 'react';
import { connect } from 'react-redux';

import { Link } from 'react-router-dom';

import { loadDevices, loadDeviceFiles, focusLocation } from '../actions';

import { FkPromisedApi } from '../api/calls';

import { Loading } from '../components/Loading';

import { generatePointDecorator } from '../common/utilities';
import MapContainer from '../components/MapContainer';

import '../../css/files.css';
import '../../css/bootstrap.min.css';

class ConcatenatedFiles extends Component {
    state = {
        details: null
    }

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
        });
    }

    render() {
        const { details } = this.state;

        if (!_.isObject(details)) {
            return <Loading />;
        }

        return (
            <div>
              <div className="">
                <table className="table table-sm">
                  <tbody>
                    {this.renderFileType('Data', 'data')}
                    {this.renderFileType('Logs', 'logs')}
                  </tbody>
                </table>
              </div>
            </div>
        );
    }

    renderFileType(name, key) {
        const { details } = this.state;
        const { device } = this.props;

        const files = details.files[key];
        const urls = device.urls[key];

        if (!_.isObject(files)) {
            return (
                <tr>
                  <th>{name}</th>
                  <td><button className="btn btn-sm btn-success" onClick={() => this.generate(urls.generate)}>Generate</button></td>
                  <td colSpan="5">None</td>
                </tr>
            );
        }

        const typeMoment = moment(files.time);
        const typeTime = typeMoment.format('MMM Do YYYY, h:mm:ss a');
        const typeAgo = typeMoment.fromNow();

        return (
            <tr>
              <th>{name}</th>
              <td><button className="btn btn-sm btn-success" onClick={() => this.generate(urls.generate)}>Generate</button></td>
              <td>{typeTime}</td>
              <td>{typeAgo}</td>
              <td>{prettyBytes(files.size)} bytes</td>
              <td><a target="_blank" rel="noopener noreferrer" href={urls.csv + "?dl=0"}>CSV</a> (<a href={urls.csv}>Download</a>)</td>
              <td><a target="_blank" rel="noopener noreferrer" href={urls.fkpb + "?dl=0"}>FKPB</a> (<a href={urls.fkpb}>Download</a>)</td>
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
        const { loadDeviceFiles, deviceId, focusLocation } = this.props;

        if (deviceId && deviceId !== prevProps.deviceId) {
            const device = this.getDevice();

            if (device.locations.entries.length > 0) {
                const center = [ ...device.locations.entries[0].coordinates, 8 ];
                focusLocation(center);
            }

            loadDeviceFiles(deviceId);
        }
    }

    renderFiles(device, deviceFiles, deviceId) {
        if (deviceFiles.all.length === 0) {
            return <Loading />;
        }

        return (
            <div className="">
              <h4>
                <div style={{ marginRight: 10, float: 'left', width: 30, height: 30, borderRadius: 10, backgroundColor: device.color }}> </div>
                <Link to={'/files'}>Back to Devices</Link>
              </h4>

              <h5>Concatenated Files</h5>

              <ConcatenatedFiles device={device} deviceId={deviceId} />

              <h5>Individual Uploads</h5>

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
        );
    }

    renderFile(file) {
        const time = moment(file.time).format('MMM Do YYYY, h:mm:ss a');

        return (
            <tr key={file.file_id}>
                <td>{file.file_id}</td>
                <td>{file.file_type_name}</td>
                <td>{prettyBytes(file.size)}</td>
                <td>{time}</td>
                <td>
                    <a target="_blank" rel="noopener noreferrer" href={file.urls.csv + "?dl=0"}>CSV</a> (<a href={file.urls.csv}>Download</a>) <span> | </span>
                    <a target="_blank" rel="noopener noreferrer" href={file.urls.json + "?dl=0"}>JSON</a> (<a href={file.urls.json}>Download</a>) <span> | </span>
                    <a href={file.urls.raw}>Raw</a>
                </td>
            </tr>
        );
    }

    renderDevice(device) {
        const fileMoment = moment(device.last_file_time);
        const time = fileMoment.format('MMM Do YYYY HH:mm');;
        const fileAgo = fileMoment.fromNow();
        const places = _.join(_(device.locations.entries).map(le => le.places).uniq().filter(s => s !== null && s !== "").value(), ", ");

        return (
            <tr className="device" key={device.device_id}>
              <td>
                    <div style={{ marginRight: 10, float: 'left', width: 20, height: 20, borderRadius: 5, backgroundColor: device.color }}> </div>
                    <Link to={'/files/' + device.device_id}>{device.device_id}</Link>
                </td>
                <td>{device.number_of_files}</td>
                <td>{prettyBytes(device.logs_size)}</td>
                <td>{prettyBytes(device.data_size)}</td>
                <td>{fileAgo} ({time})</td>
                <td>{places}</td>
            </tr>
        );
    }

    getDevice() {
        const { files, deviceId } = this.props;
        return _(files.devices).filter(d => d.device_id === deviceId).first();
    }

    renderMain() {
        const { files, deviceId } = this.props;

        if (deviceId) {
            const deviceFiles = files.filesByDevice[deviceId];
            const device = this.getDevice();

            if (!_.isObject(device) || !_.isObject(deviceFiles)) {
                return <Loading />;
            }

            return this.renderFiles(device, deviceFiles, deviceId);
        }

        return (
            <table className="table table-striped table-sm">
              <thead>
                <tr>
                  <th>ID</th>
                  <th>Files</th>
                  <th>Logs</th>
                  <th>Data</th>
                  <th>Last File</th>
                  <th>Places</th>
                </tr>
              </thead>
              <tbody>
                {files.devices.map(device => this.renderDevice(device))}
              </tbody>
            </table>
        );
    }

    render() {
        const { files } = this.props;

        if (files.devices.length === 0) {
            return <Loading />;
        }

        const features = _(files.devices)
              .map((d) => d.locations.entries.map(l => {
                  return {
                      type: 'Feature',
                      geometry: {
                          type: "Point",
                          coordinates: l.coordinates
                      },
                      properties: {
                          device: d,
                          size: 10,
                          color: d.color
                      }
                  };
              }))
              .flatten()
              .value();

        const narrowed = {
            geojson: { type: '', features: features },
            sources: [],
            focus: {
                center: files.center
            }
        };

        const pointDecorator = generatePointDecorator('constant', 'constant');

        return (
            <div className="files page">
                <div className="map-container">
                <div className="map">
                    <MapContainer style={{ height: "100%" }} containerStyle={{ width: "100%", height: "100vh" }}
                                pointDecorator={ pointDecorator } visibleFeatures={ narrowed } controls={false}
                                playbackMode={ () => false } focusFeature={ (feature) => console.log("FOCUS", feature) }
                                focusSource={ (source) => console.log("SOURCE", source) } onUserActivity={ () => false }
                                clickFeature={ (feature) => console.log("FEATURE", feature) }
                                loadMapFeatures={ () => false }
                                onChangePlaybackMode={ () => false } />
                </div>
                </div>

                <div className="main-container">
                {this.renderMain()}
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
    loadDeviceFiles,
    focusLocation
})(showWhenReady(Files, props => {
    return true;
}));
