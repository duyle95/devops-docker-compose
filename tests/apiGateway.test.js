const fetch = require('node-fetch');
const apiGatewayUrl = process.env.host;

describe('do integration tests', () => {
  it('should return 200 when fetch all container info', async () => {
    // act
    const response = await fetch(apiGatewayUrl + '/api/get-container-info');

    // assert
    expect(response.status).toBe(200);
  });
});