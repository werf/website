import logMessage from './js/chat'
import './css/style.css'
import './css/chat.css'

// Log message to console
logMessage('Its finished!!')

if (module.hot)       // eslint-disable-line no-undef
  module.hot.accept() // eslint-disable-line no-undef