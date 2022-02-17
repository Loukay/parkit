import amqp from 'amqplib';

const IMAGES_QUEUE = 'parking_images';

const createClient = async (url) => {
  const connection = await amqp.connect(url);
  const channel = await connection.createChannel();
  await channel.assertQueue(IMAGES_QUEUE);
  return channel;
};

export { createClient };
