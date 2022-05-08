function(ctx) {
  identity_id: if ctx["identity"] != null then ctx.identity.id,
  traits: ctx.identity.traits,
  state: ctx.identity.state,
  verifiable_addresses: ctx.identity.verifiable_addresses,
  request_info: {
    user_agent: ctx.request_headers["User-Agent"],
    ip: ctx.request_headers["X-Forwarded-For"],
    referer: ctx.request_headers["Referer"],
    via: ctx.request_headers["Via"]
  }
}