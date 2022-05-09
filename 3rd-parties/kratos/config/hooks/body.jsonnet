function(ctx) {
  identity_id: if ctx["identity"] != null then ctx.identity.id,
  traits: ctx.identity.traits,
  state: ctx.identity.state,
  verifiable_addresses: ctx.identity.verifiable_addresses,
}