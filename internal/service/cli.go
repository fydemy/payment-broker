package service

import (
	"fmt"
	"regexp"

	"github.com/AlecAivazis/survey/v2"
)

type CLIService interface {
	MainMenu()
	AddTenant()
	ViewTenants()
	DeleteTenant()
}

type cliService struct {
	tenantService TenantService
}

func NewCLIService(tenantService TenantService) CLIService {
	return &cliService{tenantService: tenantService}
}

func (h *cliService) MainMenu() {
	for {
		var choice string
		prompt := &survey.Select{
			Message: "Choose an action:",
			Options: []string{
				"Add Tenant",
				"View Tenants",
				"Delete Tenant",
				"Exit",
			},
		}
		survey.AskOne(prompt, &choice)

		switch choice {
		case "Add Tenant":
			h.AddTenant()
		case "View Tenants":
			h.ViewTenants()
		case "Delete Tenant":
			h.DeleteTenant()
		case "Exit":
			fmt.Println("👋 Goodbye!")
			return
		}
	}
}

func (h *cliService) AddTenant() {
	questions := []*survey.Question{
		{
			Name:     "name",
			Prompt:   &survey.Input{Message: "Enter tenant name:"},
			Validate: survey.Required,
		},
		{
			Name:     "accountID",
			Prompt:   &survey.Input{Message: "Enter Xendit Account ID:"},
			Validate: survey.Required,
		},
		{
			Name:   "webhookURL",
			Prompt: &survey.Input{Message: "Enter webhook URL:"},
			Validate: func(val interface{}) error {
				if str, ok := val.(string); !ok || str == "" {
					return fmt.Errorf("webhook URL is required")
				} else {
					// Validate URL format
					matched, _ := regexp.MatchString(`^https?://.*`, str)
					if !matched {
						return fmt.Errorf("must be a valid URL starting with http:// or https://")
					}
				}
				return nil
			},
		},
	}

	answers := struct {
		Name       string
		AccountID  string `survey:"accountID"`
		WebhookURL string `survey:"webhookURL"`
	}{}

	err := survey.Ask(questions, &answers)
	if err != nil {
		fmt.Println("❌ Error:", err)
		return
	}

	tenant, err := h.tenantService.CreateTenant(answers.Name, answers.AccountID, answers.WebhookURL)
	if err != nil {
		fmt.Println("❌ Failed to create tenant:", err)
		return
	}

	fmt.Printf("\n✅ Tenant created successfully!\n")
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	fmt.Printf("ID:          %d\n", tenant.ID)
	fmt.Printf("Name:        %s\n", tenant.Name)
	fmt.Printf("Account ID:  %s\n", tenant.AccountID)
	fmt.Printf("Webhook URL: %s\n", tenant.WebhookURL)
	fmt.Printf("API Key:     %s\n", tenant.APIKey)
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n\n")
}

func (h *cliService) ViewTenants() {
	tenants, err := h.tenantService.GetAllTenants()
	if err != nil {
		fmt.Println("❌ Error:", err)
		return
	}

	if len(tenants) == 0 {
		fmt.Println("\n📭 No tenants found\n")
		return
	}

	fmt.Println("\n🏢 Tenants List:")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	for _, tenant := range tenants {
		fmt.Printf("ID: %d | Name: %-20s | Account ID: %-24s\n", tenant.ID, tenant.Name, tenant.AccountID)
		fmt.Printf("       Webhook: %s\n", tenant.WebhookURL)
		fmt.Printf("       API Key: %s\n", tenant.APIKey)
		fmt.Println("──────────────────────────────────────────────────────────────────")
	}
	fmt.Println()
}

func (h *cliService) DeleteTenant() {
	tenants, err := h.tenantService.GetAllTenants()
	if err != nil {
		fmt.Println("❌ Error:", err)
		return
	}

	if len(tenants) == 0 {
		fmt.Println("\n📭 No tenants found\n")
		return
	}

	options := make([]string, len(tenants))
	for i, tenant := range tenants {
		options[i] = fmt.Sprintf("ID: %d - %s", tenant.ID, tenant.Name)
	}
	options = append(options, "Cancel")

	var choice string
	prompt := &survey.Select{
		Message: "Select tenant to delete:",
		Options: options,
	}
	survey.AskOne(prompt, &choice)

	if choice == "Cancel" {
		fmt.Printf("❌ Cancelled\n")
		return
	}

	var selectedID uint
	fmt.Sscanf(choice, "ID: %d", &selectedID)

	confirm := false
	confirmPrompt := &survey.Confirm{
		Message: fmt.Sprintf("Are you sure you want to delete tenant ID %d?", selectedID),
	}
	survey.AskOne(confirmPrompt, &confirm)

	if !confirm {
		fmt.Printf("❌ Cancelled\n")
		return
	}

	err = h.tenantService.DeleteTenant(selectedID)
	if err != nil {
		fmt.Println("❌ Failed to delete tenant:", err)
		return
	}

	fmt.Printf("✅ Tenant ID %d deleted successfully!\n\n", selectedID)
}
