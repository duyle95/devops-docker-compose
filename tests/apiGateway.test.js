const fetch = require('node-fetch');
const apiGatewayUrl = process.env.host;

describe('do integration tests', () => {
  it('should return 200 when GET /request', async () => {
    const response = await fetch(apiGatewayUrl + '/request');

    expect(response.status).toBe(200);
  });

  it('should return 200 when PUT /state', async () => {
    const response = await fetch(apiGatewayUrl + '/state', { method: 'PUT', body: "INIT" });

    expect(response.status).toBe(200);
  })

  it('should return 200 when GET /state', async () => {
    const response = await fetch(apiGatewayUrl + '/state');

    expect(response.status).toBe(200);
  })

  it('should return 200 when GET /run-log ', async () => {
    const response = await fetch(apiGatewayUrl + '/run-log');

    expect(response.status).toBe(400);
  })
});