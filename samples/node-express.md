Language: JavaScript/Node.js
Version: requires Node ^18
Dependencies:
  - express
  - jest

Project Structure:
└── index.js
└── package.json


## Source Files

### index.js (11 lines)
```js
const express = require('express');
const app = express();

app.get('/', (req, res) => {
  res.json({ message: 'Hello from Express!' });
});

app.listen(3000, () => {
  console.log('Server running on port 3000');
});

```

### package.json (16 lines)
```json
{
  "name": "express-sample",
  "version": "1.0.0",
  "description": "Sample Express API",
  "main": "index.js",
  "engines": {
    "node": "^18"
  },
  "dependencies": {
    "express": "^4.18.2"
  },
  "devDependencies": {
    "jest": "^29.0.0"
  }
}

```
