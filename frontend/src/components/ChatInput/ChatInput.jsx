import React, {Component} from 'react';
import './ChatInput.scss';

class ChatInput extends Component {
	render(){
		return (
			<div className = 'ChatInput'>
				<input onKeyDown={this.props.send} placeholder = "Enter a message… Hit Enter to send"/>
            </div>
        );
    };
} 

export default ChatInput;
