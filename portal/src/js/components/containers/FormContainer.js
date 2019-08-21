// @flow

import React, { Component } from 'react'

export type FormContainerProps = {
cancelText?: string;
onCancel?: () => void;

saveDisabled?: boolean;
saveText?: ?string;
onSave: () => void;

children?: any;
};

export class FormContainer extends Component {
    props: FormContainerProps;

    cancel() {
        if (this.props.onCancel) {
            this.props.onCancel();
        }
    }

    render() {
        return (
            <div className="form">
                { this.props.children }
                <div className="form-footer">
                    { this.props.onCancel !== undefined &&
                      <button className="cancel" onClick={ this.props.onCancel }>
                          { this.props.cancelText || 'Cancel' }
                      </button> }
                    <button className="save" onClick={ this.props.onSave } disabled={ this.props.saveDisabled }>
                        { this.props.saveText || 'Save' }
                    </button>
                </div>
            </div>
            );
    }
}
