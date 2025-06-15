package main

func (m MainModel) View() string {
	switch m.state {
	case WelcomeState:
		return m.welcomeModel.View()
	case StatisticsState:
		return m.statsModel.View()
	}
	return ""
}
