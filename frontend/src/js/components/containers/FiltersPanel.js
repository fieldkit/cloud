// @flow weak

import _ from 'lodash';
import moment from 'moment';
import React, { Component } from 'react';
import { Link } from 'react-router-dom';

import { API_HOST } from '../../secrets';

const panelStyle: React.CSSProperties = {
    color: "#000",
    borderTopLeftRadius: 2,
    borderTopRightRadius: 2,
    borderBottomLeftRadius: 2,
    borderBottomRightRadius: 2,
};

const containerStyle: React.CSSProperties = {
    position: 'absolute',
    width: '300px',
    zIndex: 10,
    display: 'flex',
    flexDirection: 'column',
};

const sourcesTitleStyle: React.CSSProperties = {
    padding: '10px',
    fontSize: '18px',
    fontWeight: 'bold',
    paddingBottom: '5px',
    backgroundColor: "rgba(0, 0, 0, 0.6)",
    color: 'white',
};

const listContainerStyle: React.CSSProperties = {
    backgroundColor: "rgba(255, 255, 255, 1.00)",
    padding: '10px',
};

const sourceContainerStyle: React.CSSProperties = {
    paddingBottom: '10px',
};

const sourceFirstLineStyle: React.CSSProperties = {
    fontSize: '14px',
};

class SourcePanel extends Component {
    constructor(props) {
        super(props);
        this.state = {
            expanded: false
        };
    }

    render() {
        const { info } = this.props;
        const { source, lastFeature } = info;
        const { expanded } = this.state;

        const lastFeatureDate = moment(new Date(lastFeature.properties.timestamp));
        const lastFeatureAge = lastFeatureDate.fromNow();

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
        const { onShowSource, onShowFeature, info } = this.props;
        const { source, lastFeature } = info;

        const featuresUrl = API_HOST + "/inputs/" + source.id + "/geojson?descending=true";
        const csvUrl = API_HOST + "/inputs/" + source.id + "/csv";

        return (
            <div>
                <a href={featuresUrl} target="_blank">GeoJSON</a> | <a href={csvUrl} target="_blank">CSV</a>
            </div>
        );
    }
};

export default class FiltersPanel extends Component {
    props: {
        style: React.CSSProperties
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
                    {_.values(sources).filter(r => r.source && r.summary && r.lastFeature).map((r, i) => <SourcePanel key={i} info={r} onShowFeature={onShowFeature} onShowSource={onShowSource} />)}
                </div>
            </div>
        );
    }
}
