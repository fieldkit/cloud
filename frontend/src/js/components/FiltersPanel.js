// @flow weak

import _ from 'lodash';
import moment from 'moment';
import React, { Component } from 'react';
import { Link } from 'react-router-dom';

import type { StyleSheet } from '../types';

import { API_HOST } from '../secrets';

const panelStyle: StyleSheet = {
    color: "#000",
    borderTopLeftRadius: 2,
    borderTopRightRadius: 2,
    borderBottomLeftRadius: 2,
    borderBottomRightRadius: 2,
};

const containerStyle: StyleSheet = {
    position: 'absolute',
    width: '300px',
    zIndex: 10,
    display: 'flex',
    flexDirection: 'column',
};

const sourcesTitleStyle: StyleSheet = {
    padding: '10px',
    fontSize: '18px',
    fontWeight: 'bold',
    paddingBottom: '5px',
    backgroundColor: "rgba(0, 0, 0, 0.6)",
    color: 'white',
};

const listContainerStyle: StyleSheet = {
    backgroundColor: "rgba(255, 255, 255, 1.00)",
    padding: '10px',
};

const sourceContainerStyle: StyleSheet = {
    paddingBottom: '10px',
};

const sourceFirstLineStyle: StyleSheet = {
    fontSize: '14px',
};

class SourcePanel extends Component {
    props: {
        style?: StyleSheet,
        info: any,
    }

    constructor(props) {
        super(props);
        this.state = {
            expanded: false
        };
    }

    render() {
        const { info } = this.props;
        const { source } = info;
        const { expanded } = this.state;

        const lastFeatureAge = moment(new Date(source.endTime)).fromNow();
        const numberOfFeatures = source.numberOfFeatures;

        return (
                <div style={sourceContainerStyle}>
                <div style={sourceFirstLineStyle} onClick={ () => this.setState({ expanded: !expanded }) }>
                        <span style={{ fontSize: '11px', float: 'right' }}>{numberOfFeatures} total, {lastFeatureAge}</span>
                        <Link to={"/sources/" + source.id}>
                        Device: <b>{source.name}</b>
                        </Link>
                    </div>
                    {expanded && this.renderExpanded()}
                </div>
        );
    }

    renderExpanded() {
        const { info } = this.props;
        const { source } = info;

        const featuresUrl = API_HOST + "/sources/" + source.id + "/geojson?descending=true";
        const csvUrl = API_HOST + "/sources/" + source.id + "/csv";

        return (
            <div>
                <a href={featuresUrl} target="_blank">GeoJSON</a> | <a href={csvUrl} target="_blank">CSV</a>
            </div>
        );
    }
};

export default class FiltersPanel extends Component {
    props: {
        style?: StyleSheet,
        visibleFeatures: any,
        onShowSource: () => void,
        onShowFeature: () => void,
    }

    render() {
        const { style, visibleFeatures } = this.props;
        const position = { top: 160, left: 20, bottom: 'auto', right: 'auto' };

        if (_.values(visibleFeatures.sources).length === 0) {
            return (<div></div>);
        }

        return (
            <div style={{ ...panelStyle, ...containerStyle, ...style, ...position }}>
                {this.renderSources(visibleFeatures.sources)}
            </div>
        );
    }

    renderSources(sources) {
        const { onShowSource, onShowFeature } = this.props;

        return (
            <div>
                <div style={sourcesTitleStyle}>Sources</div>
                <div style={listContainerStyle}>
                    {_.values(sources).filter(r => r.source && r.summary).map((r, i) => <SourcePanel key={i} info={r} onShowFeature={onShowFeature} onShowSource={onShowSource} />)}
                </div>
            </div>
        );
    }
}
