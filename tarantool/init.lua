box.cfg {
  listen = 3301
}

box.once("schema", function()
  box.schema.space.create("kvstore")
  box.space.kvstore:format({
    { name = "key",   type = "string" },
    { name = "value", type = "string" }
  })
  box.space.kvstore:create_index("primary", {
    parts = { "key" }
  })
end)
