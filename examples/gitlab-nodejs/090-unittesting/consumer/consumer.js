const amqp = require('amqplib')
const mailgun = require("mailgun-js");

const sendEmail = async (email) => {
  try {
    const mg = mailgun({apiKey: process.env.MAILGUN_APIKEY, domain: process.env.MAILGUN_DOMAIN, host: "api.eu.mailgun.net"});
    if (email != null) {
      email.from = "Mailgun Sandbox <postmaster@" + process.env.MAILGUN_DOMAIN + ">"
      email.subject = "Welcome to Krovatka!"
      email.html = "<p>" + email.text + "</p>"
      console.log("Email: " + JSON.stringify(email));
    }
    const sent = await mg.messages().send(email);
    return sent;
  } catch (error) {
    console.error(error)
  }
}

(async () => {
  try {
    const conn = await amqp.connect(process.env.AMQP_URI + '?heartbeat=5s')
    const ch = await conn.createChannel()
    const queueName = 'onboard_emails'
    // create queue
    await ch.assertQueue(queueName, { durable: true })
    // subscribe
    console.log('Started to listen')
    await ch.consume(queueName, msg => {
      try {
        let email = null;
        // parse
        try {
          email = JSON.parse(msg.content.toString());
        } catch (error) {
          console.error(error)
        }
      const res = sendEmail(email)
      } catch (error) {
        // not accepting the message, because we have failed sending email
        ch.nack(msg);
        console.error(error)
      }
    });
  } catch (error) {
      console.error(error)
  }
})()
