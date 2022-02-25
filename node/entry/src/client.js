import { createClient as createRedisClient } from 'redis';

const createClient = async ({ url } = {}) => {

    const client = createRedisClient({ url });

    client.connect();

    return client;

}

const Client = await createClient({
    url: process.env.REDIS_URL
});

export default Client;
