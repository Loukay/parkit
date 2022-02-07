from flask import Flask, render_template, request, make_response
import cv2
import numpy as np 
import torch
import torch.nn.functional as F
import json



normalize = lambda x, mean=0.5, std=0.25: (x - mean) / std
denormalize = lambda x, mean=0.5, std=0.25: x * std + mean
resize = torch.nn.Upsample(size=(128, 128), mode='bilinear', align_corners=False)

model = None
spots = None

app = Flask(__name__)

@app.route('/predict', methods=['POST'])
def predict():
    img = np.frombuffer(request.data, dtype=np.uint8)
    results = get_prediction(img,model)
    return 'nadi'

def order_points(pts):

	rect = np.zeros((4, 2), dtype = "float32")

	s = pts.sum(axis = 1)
	rect[0] = pts[np.argmin(s)]
	rect[2] = pts[np.argmax(s)]

	diff = np.diff(pts, axis = 1)
	rect[1] = pts[np.argmin(diff)]
	rect[3] = pts[np.argmax(diff)]
	return rect

def four_point_transform(image, pts):

    rect = order_points(pts)
    (tl, tr, br, bl) = rect

    widthA = np.sqrt(((br[0] - bl[0]) ** 2) + ((br[1] - bl[1]) ** 2))
    widthB = np.sqrt(((tr[0] - tl[0]) ** 2) + ((tr[1] - tl[1]) ** 2))
    maxWidth = max(int(widthA), int(widthB))

    heightA = np.sqrt(((tr[0] - br[0]) ** 2) + ((tr[1] - br[1]) ** 2))
    heightB = np.sqrt(((tl[0] - bl[0]) ** 2) + ((tl[1] - bl[1]) ** 2))
    maxHeight = max(int(heightA), int(heightB))

    dst = np.array([
        [0, 0],
        [maxWidth - 1, 0],
        [maxWidth - 1, maxHeight - 1],
        [0, maxHeight - 1]], dtype = "float32")

    M = cv2.getPerspectiveTransform(rect, dst)
    warped = cv2.warpPerspective(image, M, (maxWidth, maxHeight))

    return warped

def get_coordinates(polygon):
    coordinates = []
    for c in polygon:
        coordinates.append([c[0], c[1]])
    return coordinates



def get_prediction(img_data,model):
    # Image
    
    image = cv2.imdecode(img_data, cv2.IMREAD_COLOR)[..., ::-1]

    for i, p in enumerate(spots):
        pts = np.array(p, dtype=np.int64)
        warped = four_point_transform(image, pts)
        im = cv2.resize(warped, (128,128))
        im = np.ascontiguousarray(np.asarray(im).transpose((2, 0, 1)))  # HWC to CHW
        im = torch.tensor(im).float().unsqueeze(0) / 255.0  # to Tensor, to BCWH, rescale
        im = resize(normalize(im))
        results = model(im)
        p = F.softmax(results, dim=1)  # probabilities
        i = p.argmax()  # max index
        #print(f'{file} prediction: {i} ({p[0, i]:.2f})')
        print(i)


def classify(model, size=128, file='../datasets/mnist/test/3/30.png', plot=False):
    # YOLOv5 classification model inference


    resize = torch.nn.Upsample(size=(size, size), mode='bilinear', align_corners=False)  # image resize

    # Image
    im = cv2.imread(str(file))[..., ::-1]  # HWC, BGR to RGB
    im = np.ascontiguousarray(np.asarray(im).transpose((2, 0, 1)))  # HWC to CHW
    im = torch.tensor(im).float().unsqueeze(0) / 255.0  # to Tensor, to BCWH, rescale
    im = resize(normalize(im))

    # Inference
    results = model(im)
    p = F.softmax(results, dim=1)  # probabilities
    i = p.argmax()  # max index
    #print(f'{file} prediction: {i} ({p[0, i]:.2f})')

    return i

if __name__ == '__main__':

    polygons = open('./spots.json')
    spots = json.load(polygons)
    polygons.close() 

    model = torch.load('./models/model.pt', map_location=torch.device('cuda'))['model'].float() #the keys in the listOfKeys

    print('detector running')

    app.run(host='0.0.0.0')


