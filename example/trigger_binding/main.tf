provider "auth0" {}

resource "auth0_action" "action_foo" {
    name = "Test Trigger Binding Foo {{.random}}"
    supported_triggers {
        id = "post-login"
        version = "v2"
    }
    code = <<-EOT
	exports.onContinuePostLogin = async (event, api) => {
		console.log("foo")
	};"
	EOT
    deploy = true
}

resource "auth0_action" "action_bar" {
    name = "Test Trigger Binding Bar {{.random}}"
    supported_triggers {
        id = "post-login"
        version = "v2"
    }
    code = <<-EOT
	exports.onContinuePostLogin = async (event, api) => {
		console.log("bar")
	};"
	EOT
    deploy = true
}

resource "auth0_trigger_binding" "login_flow" {
    trigger = "post-login"
    actions {
        id = auth0_action.action_foo.id
        display_name = auth0_action.action_foo.name
    }
    actions {
        id = auth0_action.action_bar.id
        display_name = auth0_action.action_bar.name
    }
}
