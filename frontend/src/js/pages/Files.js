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

class DeviceNotes extends Component {
    save(ev) {
        const { device } = this.props;

        ev.preventDefault();

        const update = {
            deviceId: device.device_id,
            name: (ev.target.name.value !== "") ? ev.target.name.value : null,
            notes: (ev.target.notes.value !== "") ? ev.target.notes.value : null,
        };

        if (update.name == null && update.notes == null) {
            return;
        }

        fetch(device.urls.details, {
            method: 'POST',
            body: JSON.stringify(update)
        }).then((res) => {
            return res.json();
        }).then(data =>{
            console.log(data);
        });
    }

    renderNote(device, index, note) {
        const noteMoment = moment(note.time);
        const noteAgo = noteMoment.fromNow();

        return (
            <tr key={index}>
                <th> {noteAgo} </th>
                <td> {note.name && "Name changed to " + note.name} {note.notes && "Notes: " + note.notes} </td>
            </tr>
        );
    }

    render() {
        const { device, details } = this.props;

        if (!_.isObject(details)) {
            return <Loading />;
        }

        return (
            <div>
                <table className="notes">
                    <tbody>
                    {details.notes.map((note, idx) => {
                        return this.renderNote(device, idx, note);
                    })}
                    </tbody>
                </table>
                <form className="" onSubmit={(ev) => this.save(ev)}>
                    <div className="form-row">
                        <div className="form-group col-md-6">
                            <label htmlFor="name">Name</label>
                            <input type="text" name="name" className="form-control" placeholder="Name" />
                        </div>
                        <div className="form-group col-md-6">
                            <label htmlFor="notes">Notes</label>
                            <input type="text" name="notes" className="form-control" placeholder="Notes" />
                        </div>
                    </div>
                    <div className="form-row">
                        <div className="col-md-6">
                        <button type="submit" className="btn btn-primary">Save</button>
                        </div>
                    </div>
                </form>
            </div>
        );
    }
}

class ConcatenatedFiles extends Component {
    generate(url) {
        fetch(url, {}).then((res) => {
            return res.json();
        }).then(data =>{
            console.log(data);
        });
    }

    render() {
        const { details } = this.props;

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
        const { device, details } = this.props;

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
                <td><a href={urls.fkpb}>FKPB</a></td>
            </tr>
        );
    }
}

class Files extends Component {
    props: {
        loadDevices: any => void,
    };

    state = {
        details: null
    };

    refreshDevice(deviceId) {
        const { loadDeviceFiles } = this.props;

        loadDeviceFiles(deviceId);

        FkPromisedApi.queryDeviceFilesDetails(deviceId).then(details => {
            this.setState({
                details
            });
        });
    }

    componentWillMount() {
        const { loadDevices, deviceId } = this.props;

        loadDevices();

        if (deviceId) {
            this.refreshDevice(deviceId);
        }
    }

    componentDidUpdate(prevProps, prevState, snapshot) {
        const { deviceId, focusLocation } = this.props;

        if (deviceId && deviceId !== prevProps.deviceId) {
            const device = this.getDevice();

            if (device.locations.entries.length > 0) {
                const center = [ ...device.locations.entries[0].coordinates, 8 ];
                focusLocation(center);
            }

            this.refreshDevice(deviceId);
        }
    }

    renderFiles(device, deviceFiles, deviceId) {
        const { details } = this.state;

        if (deviceFiles.all.length === 0) {
            return <Loading />;
        }

        return (
            <div className="">
                <h4>
                    <div style={{ marginRight: 10, float: 'left', width: 30, height: 30, borderRadius: 10, backgroundColor: device.color }}> </div>
                    {deviceId} {device.name}
                </h4>

                <div style={{ fontSize: 16, fontWeight: 'bold', marginBottom: 10 }}>
                    <Link to={'/files'}>Back to Devices</Link>
                </div>

                <h5>Notes</h5>

                <DeviceNotes device={device} details={details} />

                <h5>Concatenated Files</h5>

                <ConcatenatedFiles device={device} deviceId={deviceId} details={details} />

                <h5>Individual Files</h5>

                <div className="">
                    <table className="table table-striped table-sm">
                        <thead>
                            <tr>
                                <th>File</th>
                                <th>Type</th>
                                <th>Size</th>
                                <th>Uploaded</th>
                                <th></th>
                                <th></th>
                                <th></th>
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
                <td>{file.corrupted ? "Corrupt" : ""}</td>
                <td> <a target="_blank" rel="noopener noreferrer" href={file.urls.csv + "?dl=0"}>CSV</a> (<a href={file.urls.csv}>Download</a>) </td>
                <td> <a href={file.urls.raw}>FKPB</a> </td>
            </tr>
        );
    }

    renderDevice(device) {
        const fileMoment = moment(device.last_file_time);
        // const time = fileMoment.format('MMM Do YYYY HH:mm');
        const fileAgo = fileMoment.fromNow();
        const places = _.join(_(device.locations.entries).map(le => le.places).uniq().filter(s => s !== null && s !== "").value(), ", ");

        return (
            <tr className="device" key={device.device_id}>
                <td>
                    <div style={{ marginRight: 10, float: 'left', width: 20, height: 20, borderRadius: 5, backgroundColor: device.color }}> </div>
                    <Link to={'/files/' + device.device_id}>{device.device_id}</Link>
                </td>
                <td>{device.name}</td>
                <td>{device.number_of_files}</td>
                <td>{prettyBytes(device.logs_size)}</td>
                <td>{prettyBytes(device.data_size)}</td>
                <td>{fileAgo}</td>
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
                        <th>Name</th>
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
