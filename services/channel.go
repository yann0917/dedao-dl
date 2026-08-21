package services

// ChannelInfo 获取学习圈频道信息
func (s *Service) ChannelInfo(channelID int) (info *ChannelInfo, err error) {

	body, err := s.reqChannelInfo(channelID)
	if err != nil {
		return
	}
	defer body.Close()

	if err = handleJSONParse(body, &info); err != nil {
		return
	}
	return
}

// ChannelHomepage 获取学习圈频道首页分类
func (s *Service) ChannelHomepage(channelID int) (cats []ChannelHomepageCategory, err error) {

	body, err := s.reqChannelHomepage(channelID)
	if err != nil {
		return
	}
	defer body.Close()

	if err = handleJSONParse(body, &cats); err != nil {
		return
	}
	return
}

// ChannelVipInfo 获取学习圈VIP/权限信息
func (s *Service) ChannelVipInfo(channelID int) (info *ChannelVipInfo, err error) {
	body, err := s.reqChannelVipInfo(channelID)
	if err != nil {
		return
	}
	defer body.Close()

	if err = handleJSONParse(body, &info); err != nil {
		return
	}
	return
}

// ChannelTopicDetail 根据主题ID获取主题详情列表
// 接口返回 {"topic": {...}} 包装结构，这里解包后直接返回主题对象
func (s *Service) ChannelTopicDetail(topicID int) (topic *ChannelTopicCategory, err error) {
	body, err := s.reqChannelTopicDetail(topicID)
	if err != nil {
		return
	}
	defer body.Close()

	var detail *ChannelTopicDetail
	if err = handleJSONParse(body, &detail); err != nil {
		return
	}
	if detail == nil {
		return nil, nil
	}
	return detail.Topic, nil
}
