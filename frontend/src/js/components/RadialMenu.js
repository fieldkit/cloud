import ReactRadial from 'react-radial';

export class RadialMenu extends ReactRadial {
    componentWillMount() {
        const { position } = this.props;

        this.setState({
            enabled: true,
            cx: position.x,
            cy: position.y,
            data: this._getDataObject(0, 0, { ...this.props })
        });
    }

    componentWillUpdate(nextProps, nextState) {
        const { enabled: oldEnabled } = this.state;
        const { enabled: newEnabled } = nextState;

        if (oldEnabled && !newEnabled) {
            this.props.onClosed();
        }
    }

    componentWillUnmount() {
    }

    componentWillReceiveProps(nextProps) {
        const { buttons: oldButtons } = this.props;
        const { buttons: newButtons } = nextProps;
        const { enabled } = this.state;

        /*
        if (enabled) {
            if (oldButtons != newButtons) {
                console.log("BUTTONS CHANGED");

                const { position } = this.props;

                const data = this._getDataObject(this.props.innerRadius, this.props.outerRadius, { ...nextProps });

                this.setState({
                    enabled: true,
                    cx: position.x,
                    cy: position.y,
                    data: data
                });
            }
        }
        */
    }

    render() {
        return super.render();
    }
}

