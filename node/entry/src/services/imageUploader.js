import Client from '../client.js';

import Logger from '../logger.js';

const uploadImage = (req, res) => {
    res.send();
    const image = req.body;
    // TODO: Make sure it's an image
    Client.publish('parking_images', image)
    // TODO: Get parking data from the request and include the ID in the log
    Logger.log('info', 'Photo uploaded from #23564');
}

export { uploadImage };