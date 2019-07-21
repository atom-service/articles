package provider

import (
	"context"

	"github.com/joho/godotenv"
	"github.com/yinxulai/grpc-module-articles/models"
	"github.com/yinxulai/grpc-module-articles/standard"
)

// NewService NewService
func NewService() *Service {
	godotenv.Load()
	service := new(Service)
	return service
}

// Service Service
type Service struct {
}

// CreateArticle Create
func (srv *Service) CreateArticle(ctx context.Context, req *standard.CreateArticleRequest) (resp *standard.CreateArticleResponse, err error) {
	var count uint64
	resp = new(standard.CreateArticleResponse)
	// TODO: 检查分类是否存在
	if req.OwnerCategory == 0 {
		resp.State = standard.State_PARAMS_INVALID
		resp.Message = "无效的分类"
		return resp, nil
	}

	err = countCategoryByIDNamedStmt.GetContext(ctx, &count, map[string]interface{}{"ID": req.OwnerCategory})
	if err != nil {
		resp.State = standard.State_DB_OPERATION_FATLURE
		resp.Message = err.Error()
		return resp, nil
	}

	if count <= 0 {
		resp.State = standard.State_ARTICLE_NOT_EXIST
		resp.Message = "该分类不存在"
		return resp, nil
	}

	// 执行
	_, err = insertArticleNamedStmt.ExecContext(ctx, req)
	if err != nil {
		resp.State = standard.State_DB_OPERATION_FATLURE
		resp.Message = err.Error()
		return resp, nil
	}

	resp.State = standard.State_SUCCESS
	resp.Message = "创建成功"

	return resp, nil
}

// QueryArticleByID QueryArticleByID
func (srv *Service) QueryArticleByID(ctx context.Context, req *standard.QueryArticleByIDRequest) (resp *standard.QueryArticleByIDResponse, err error) {
	articles := []*models.Article{}
	resp = new(standard.QueryArticleByIDResponse)

	if req.ID == 0 {
		resp.State = standard.State_PARAMS_INVALID
		resp.Message = "无效的 ID"
		return resp, nil
	}

	rows, err := queryArticleByIDNamedStmt.QueryxContext(ctx, req)
	if err != nil {
		resp.State = standard.State_DB_OPERATION_FATLURE
		resp.Message = err.Error()
		return resp, nil
	}

	for rows.Next() {
		var localArticle models.Article
		err = rows.StructScan(&localArticle)
		if err == nil {
			articles = append(articles, &localArticle)
		}
	}

	if len(articles) <= 0 { // 没有找到用户
		resp.State = standard.State_ARTICLE_NOT_EXIST
		resp.Message = "该文章不存在"
		return resp, nil
	}

	resp.State = standard.State_SUCCESS
	resp.Data = articles[0].OutProtoStruct()
	resp.Message = "查询成功"
	return resp, nil
}

// DeleteArticleByID DeleteArticleByID
func (srv *Service) DeleteArticleByID(ctx context.Context, req *standard.DeleteArticleByIDRequest) (resp *standard.DeleteArticleByIDResponse, err error) {
	var count uint64
	resp = new(standard.DeleteArticleByIDResponse)

	if req.ID == 0 {
		resp.State = standard.State_PARAMS_INVALID
		resp.Message = "无效的 ID"
		return resp, nil
	}

	err = countArticleByIDNamedStmt.GetContext(ctx, &count, req)
	if err != nil {
		resp.State = standard.State_DB_OPERATION_FATLURE
		resp.Message = err.Error()
		return resp, nil
	}

	if count <= 0 {
		resp.State = standard.State_ARTICLE_NOT_EXIST
		resp.Message = "该文章不存在"
		return resp, nil
	}

	_, err = deleteArticleByIDNamedStmt.ExecContext(ctx, req)
	if err != nil {
		resp.State = standard.State_DB_OPERATION_FATLURE
		resp.Message = err.Error()
		return resp, nil
	}

	resp.State = standard.State_SUCCESS
	resp.Message = "删除成功"

	return resp, nil
}

// UpdateArticleByID UpdateArticleByID
func (srv *Service) UpdateArticleByID(ctx context.Context, req *standard.UpdateArticleByIDRequest) (resp *standard.UpdateArticleByIDResponse, err error) {
	// 检查是否存在该记录
	var count uint64
	resp = new(standard.UpdateArticleByIDResponse)

	if req.ID == 0 {
		resp.State = standard.State_PARAMS_INVALID
		resp.Message = "无效的 ID"
		return resp, nil
	}

	err = countArticleByIDNamedStmt.GetContext(ctx, &count, req)
	if err != nil {
		resp.State = standard.State_DB_OPERATION_FATLURE
		resp.Message = err.Error()
		return resp, nil
	}

	if count <= 0 {
		resp.State = standard.State_ARTICLE_NOT_EXIST
		resp.Message = "该文章不存在"
		return resp, nil
	}

	req.Data.ID = req.ID
	_, err = updateArticleByIDNamedStmt.ExecContext(ctx, req.Data)
	if err != nil {
		resp.State = standard.State_DB_OPERATION_FATLURE
		resp.Message = err.Error()
		return resp, nil
	}

	resp.State = standard.State_SUCCESS
	resp.Message = "更新成功"

	return resp, nil
}

// QueryArticleByOwner QueryArticleByOwner
func (srv *Service) QueryArticleByOwner(ctx context.Context, req *standard.QueryArticleByOwnerRequest) (resp *standard.QueryArticleByOwnerResponse, err error) {
	var count uint64
	articles := []*models.Article{}
	stdarticles := []*standard.Article{}
	resp = new(standard.QueryArticleByOwnerResponse)
	// 查询 Owner 文章是否存在

	if req.Owner == 0 {
		resp.State = standard.State_PARAMS_INVALID
		resp.Message = "无效的 Owner"
		return resp, nil
	}

	// 查询记录总数
	err = countArticleByOwnerNamedStmt.GetContext(ctx, &count, req)
	if err != nil {
		resp.State = standard.State_DB_OPERATION_FATLURE
		resp.Message = err.Error()
		return resp, nil
	}

	rows, err := queryArticleByOwnerNamedStmt.QueryxContext(ctx, req)
	if err != nil {
		resp.State = standard.State_DB_OPERATION_FATLURE
		resp.Message = err.Error()
		return resp, nil
	}

	for rows.Next() {
		var localArticle models.Article
		err = rows.StructScan(&localArticle)
		if err == nil {
			articles = append(articles, &localArticle)
		}
	}

	for _, article := range articles {
		stdarticles = append(stdarticles, article.OutProtoStruct())
	}

	resp.State = standard.State_SUCCESS
	resp.Data = stdarticles
	resp.Total = count
	resp.Message = "查询成功"

	return resp, nil
}

// QueryArticleByOwnerCategory QueryByOwnerCategory
func (srv *Service) QueryArticleByOwnerCategory(ctx context.Context, req *standard.QueryArticleByOwnerCategoryRequest) (resp *standard.QueryArticleByOwnerCategoryResponse, err error) {
	// 检查是否存在该记录
	var count uint64
	articles := []*models.Article{}
	stdarticles := []*standard.Article{}
	resp = new(standard.QueryArticleByOwnerCategoryResponse)

	// 简单检查一下分类
	if req.OwnerCategory == 0 {
		resp.State = standard.State_PARAMS_INVALID
		resp.Message = "无效的分类"
		return resp, nil
	}

	// 查询 Category 是否存在
	err = countCategoryByIDNamedStmt.GetContext(ctx, &count, map[string]interface{}{"ID": req.OwnerCategory})
	if err != nil {
		resp.State = standard.State_DB_OPERATION_FATLURE
		resp.Message = err.Error()
		return resp, nil
	}

	if count <= 0 {
		resp.State = standard.State_ARTICLE_NOT_EXIST
		resp.Message = "该分类不存在"
		return resp, nil
	}

	// 查询记录总数
	err = countArticleByOwnerCategoryNamedStmt.GetContext(ctx, &count, req)
	if err != nil {
		resp.State = standard.State_DB_OPERATION_FATLURE
		resp.Message = err.Error()
		return resp, nil
	}

	// 执行条件查询
	rows, err := queryArticleByOwnerCategoryNamedStmt.QueryxContext(ctx, req)
	if err != nil {
		resp.State = standard.State_DB_OPERATION_FATLURE
		resp.Message = err.Error()
		return resp, nil
	}

	for rows.Next() {
		var localArticle models.Article
		err = rows.StructScan(&localArticle)
		if err == nil {
			articles = append(articles, &localArticle)
		}
	}

	for _, article := range articles {
		stdarticles = append(stdarticles, article.OutProtoStruct())
	}

	resp.State = standard.State_SUCCESS
	resp.Data = stdarticles
	resp.Total = count
	resp.Message = "查询成功"
	return resp, nil
}

// QueryLabelByID QueryLabelByID
func (srv *Service) QueryLabelByID(ctx context.Context, req *standard.QueryLabelByIDRequest) (resp *standard.QueryLabelByIDResponse, err error) {
	labels := []*models.Label{}
	resp = new(standard.QueryLabelByIDResponse)

	rows, err := queryLabelByIDNamedStmt.QueryxContext(ctx, req)
	if err != nil {
		resp.State = standard.State_DB_OPERATION_FATLURE
		resp.Message = err.Error()
		return resp, nil
	}

	for rows.Next() {
		var localLabel models.Label
		err = rows.StructScan(&localLabel)
		if err == nil {
			labels = append(labels, &localLabel)
		}
	}

	if len(labels) <= 0 {
		resp.State = standard.State_LABEL_NOT_EXIST
		resp.Message = "该标签不存在"
		return resp, nil
	}

	resp.State = standard.State_SUCCESS
	resp.Data = labels[0].OutProtoStruct()
	resp.Message = "查询成功"

	return resp, nil
}

// UpdateLabelByID UpdateLabelByID
func (srv *Service) UpdateLabelByID(ctx context.Context, req *standard.UpdateLabelByIDRequest) (resp *standard.UpdateLabelByIDResponse, err error) {
	// 检查是否存在该记录
	var count uint64
	resp = new(standard.UpdateLabelByIDResponse)

	if req.ID == 0 {
		resp.State = standard.State_PARAMS_INVALID
		resp.Message = "无效的 ID"
		return resp, nil
	}

	err = countLabelByIDNamedStmt.GetContext(ctx, &count, req)
	if err != nil {
		resp.State = standard.State_DB_OPERATION_FATLURE
		resp.Message = err.Error()
		return resp, nil
	}

	if count <= 0 {
		resp.State = standard.State_LABEL_NOT_EXIST
		resp.Message = "该标签不存在"
		return resp, nil
	}

	req.Data.ID = req.ID
	_, err = updateLabelByIDNamedStmt.ExecContext(ctx, req.Data)
	if err != nil {
		resp.State = standard.State_DB_OPERATION_FATLURE
		resp.Message = err.Error()
		return resp, nil
	}

	resp.State = standard.State_SUCCESS
	resp.Message = "更新成功"

	return resp, nil
}

// DeleteLabelByID DeleteLabelByID
func (srv *Service) DeleteLabelByID(ctx context.Context, req *standard.DeleteLabelByIDRequest) (resp *standard.DeleteLabelByIDResponse, err error) {
	// 检查是否存在该记录
	var count uint64
	resp = new(standard.DeleteLabelByIDResponse)

	if req.ID == 0 {
		resp.State = standard.State_PARAMS_INVALID
		resp.Message = "无效的 ID"
		return resp, nil
	}

	err = countLabelByIDNamedStmt.GetContext(ctx, &count, req)
	if err != nil {
		resp.State = standard.State_DB_OPERATION_FATLURE
		resp.Message = err.Error()
		return resp, nil
	}

	if count <= 0 {
		resp.State = standard.State_LABEL_NOT_EXIST
		resp.Message = "该标签不存在"
		return resp, nil
	}

	_, err = deleteLabelByIDNamedStmt.ExecContext(ctx, req)
	if err != nil {
		resp.State = standard.State_DB_OPERATION_FATLURE
		resp.Message = err.Error()
		return resp, nil
	}

	resp.State = standard.State_SUCCESS
	resp.Message = "删除成功"

	return resp, nil
}

// CreateLabelByOwner CreateLabelByOwner
func (srv *Service) CreateLabelByOwner(ctx context.Context, req *standard.CreateLabelByOwnerRequest) (resp *standard.CreateLabelByOwnerResponse, err error) {
	// 检查是否存在该记录
	var count uint64
	resp = new(standard.CreateLabelByOwnerResponse)

	if req.Owner == 0 {
		resp.State = standard.State_PARAMS_INVALID
		resp.Message = "无效的文章 ID"
		return resp, nil
	}

	err = countArticleByIDNamedStmt.GetContext(ctx, &count, map[string]interface{}{"ID": req.Owner})
	if err != nil {
		resp.State = standard.State_DB_OPERATION_FATLURE
		resp.Message = err.Error()
		return resp, nil
	}

	if count <= 0 {
		resp.State = standard.State_ARTICLE_NOT_EXIST
		resp.Message = "该文章不存在"
		return resp, nil
	}

	req.Label.Owner = req.Owner
	_, err = insertLabelNamedStmt.ExecContext(ctx, req.Label)
	if err != nil {
		resp.State = standard.State_DB_OPERATION_FATLURE
		resp.Message = err.Error()
		return resp, nil
	}

	resp.State = standard.State_SUCCESS
	resp.Message = "创建成功"

	return resp, nil
}

// QueryLabelByOwner QueryLabelByOwner
func (srv *Service) QueryLabelByOwner(ctx context.Context, req *standard.QueryLabelByOwnerRequest) (resp *standard.QueryLabelByOwnerResponse, err error) {
	var count uint64
	labels := []*models.Label{}
	stdlabels := []*standard.Label{}

	resp = new(standard.QueryLabelByOwnerResponse)
	// 查询 Owner 文章是否存在

	if req.Owner == 0 {
		resp.State = standard.State_PARAMS_INVALID
		resp.Message = "无效的文章 ID"
		return resp, nil
	}

	err = countArticleByIDNamedStmt.GetContext(ctx, &count, map[string]interface{}{"ID": req.Owner})
	if err != nil {
		resp.State = standard.State_DB_OPERATION_FATLURE
		resp.Message = err.Error()
		return resp, nil
	}

	if count <= 0 {
		resp.State = standard.State_ARTICLE_NOT_EXIST
		resp.Message = "该文章不存在"
		return resp, nil
	}

	// 查询记录总数
	err = countLabelByOwnerNamedStmt.GetContext(ctx, &count, req)
	if err != nil {
		resp.State = standard.State_DB_OPERATION_FATLURE
		resp.Message = err.Error()
		return resp, nil
	}

	rows, err := queryLabelByOwnerNamedStmt.QueryxContext(ctx, req)
	if err != nil {
		resp.State = standard.State_DB_OPERATION_FATLURE
		resp.Message = err.Error()
		return resp, nil
	}

	for rows.Next() {
		var localLabel models.Label
		err = rows.StructScan(&localLabel)
		if err == nil {
			labels = append(labels, &localLabel)
		}
	}

	for _, label := range labels {
		stdlabels = append(stdlabels, label.OutProtoStruct())
	}

	resp.State = standard.State_SUCCESS
	resp.Data = stdlabels
	resp.Total = count
	resp.Message = "查询成功"

	return resp, nil
}

// CreateCategory CreateCategory
func (srv *Service) CreateCategory(ctx context.Context, req *standard.CreateCategoryRequest) (resp *standard.CreateCategoryResponse, err error) {
	resp = new(standard.CreateCategoryResponse)
	// 执行
	_, err = insertCategoryNamedStmt.ExecContext(ctx, req)
	if err != nil {
		resp.State = standard.State_DB_OPERATION_FATLURE
		resp.Message = err.Error()
		return resp, nil
	}

	resp.State = standard.State_SUCCESS
	resp.Message = "创建成功"

	return resp, nil
}

// UpdateCategoryByID UpdateCategoryByID
func (srv *Service) UpdateCategoryByID(ctx context.Context, req *standard.UpdateCategoryByIDRequest) (resp *standard.UpdateCategoryByIDResponse, err error) {
	// 检查是否存在该记录
	var count uint64
	resp = new(standard.UpdateCategoryByIDResponse)

	if req.ID == 0 {
		resp.State = standard.State_PARAMS_INVALID
		resp.Message = "无效的 ID"
		return resp, nil
	}

	err = countCategoryByIDNamedStmt.GetContext(ctx, &count, req)
	if err != nil {
		resp.State = standard.State_DB_OPERATION_FATLURE
		resp.Message = err.Error()
		return resp, nil
	}

	if count <= 0 {
		resp.State = standard.State_CATEGORY_NOT_EXIST
		resp.Message = "该分类不存在"
		return resp, nil
	}

	req.Data.ID = req.ID
	_, err = updateCategoryByIDNamedStmt.ExecContext(ctx, req.Data)
	if err != nil {
		resp.State = standard.State_DB_OPERATION_FATLURE
		resp.Message = err.Error()
		return resp, nil
	}

	resp.State = standard.State_SUCCESS
	resp.Message = "更新成功"

	return resp, nil
}

// DeleteCategoryByID DeleteCategoryByID
func (srv *Service) DeleteCategoryByID(ctx context.Context, req *standard.DeleteCategoryByIDRequest) (resp *standard.DeleteCategoryByIDResponse, err error) {
	// 检查是否存在该记录
	var count uint64
	resp = new(standard.DeleteCategoryByIDResponse)

	if req.ID == 0 {
		resp.State = standard.State_PARAMS_INVALID
		resp.Message = "无效的 ID"
		return resp, nil
	}

	err = countCategoryByIDNamedStmt.GetContext(ctx, &count, req)
	if err != nil {
		resp.State = standard.State_DB_OPERATION_FATLURE
		resp.Message = err.Error()
		return resp, nil
	}

	if count <= 0 {
		resp.State = standard.State_CATEGORY_NOT_EXIST
		resp.Message = "该分类不存在"
		return resp, nil
	}

	// TODO: 检查 该分类 下的文章以及子分类
	_, err = deleteCategoryByIDNamedStmt.ExecContext(ctx, req)
	if err != nil {
		resp.State = standard.State_DB_OPERATION_FATLURE
		resp.Message = err.Error()
		return resp, nil
	}

	resp.State = standard.State_SUCCESS
	resp.Message = "删除成功"

	return resp, nil
}

// QueryCategoryByOwner QueryCategoryByOwner
func (srv *Service) QueryCategoryByOwner(ctx context.Context, req *standard.QueryCategoryByOwnerRequest) (resp *standard.QueryCategoryByOwnerResponse, err error) {
	var count uint64
	categorys := []*models.Category{}
	stdcategorys := []*standard.Category{}
	resp = new(standard.QueryCategoryByOwnerResponse)

	err = countCategoryByOwnerNamedStmt.GetContext(ctx, &count, req)
	if err != nil {
		resp.State = standard.State_DB_OPERATION_FATLURE
		resp.Message = err.Error()
		return resp, nil
	}

	rows, err := queryCategoryByOwnerNamedStmt.QueryxContext(ctx, req)
	if err != nil {
		resp.State = standard.State_DB_OPERATION_FATLURE
		resp.Message = err.Error()
		return resp, nil
	}

	for rows.Next() {
		var localCategory models.Category
		err = rows.StructScan(&localCategory)
		if err == nil {
			categorys = append(categorys, &localCategory)
		}
	}

	for _, category := range categorys {
		stdcategorys = append(stdcategorys, category.OutProtoStruct())
	}

	resp.State = standard.State_SUCCESS
	resp.Data = stdcategorys
	resp.Total = count
	resp.Message = "查询成功"

	return resp, nil
}

// QueryCategoryByOwnerCategory QueryCategoryByOwnerCategory
func (srv *Service) QueryCategoryByOwnerCategory(ctx context.Context, req *standard.QueryCategoryByOwnerCategoryRequest) (resp *standard.QueryCategoryByOwnerCategoryResponse, err error) {
	var count uint64
	categorys := []*models.Category{}
	stdcategorys := []*standard.Category{}
	resp = new(standard.QueryCategoryByOwnerCategoryResponse)

	// TODO: 查询父分类是否存在
	err = countCategoryByOwnerCategoryNamedStmt.GetContext(ctx, &count, req)
	if err != nil {
		resp.State = standard.State_DB_OPERATION_FATLURE
		resp.Message = err.Error()
		return resp, nil
	}

	rows, err := queryCategoryByOwnerCategoryNamedStmt.QueryxContext(ctx, req)
	if err != nil {
		resp.State = standard.State_DB_OPERATION_FATLURE
		resp.Message = err.Error()
		return resp, nil
	}

	for rows.Next() {
		var localCategory models.Category
		err = rows.StructScan(&localCategory)
		if err == nil {
			categorys = append(categorys, &localCategory)
		}
	}

	for _, category := range categorys {
		stdcategorys = append(stdcategorys, category.OutProtoStruct())
	}

	resp.State = standard.State_SUCCESS
	resp.Data = stdcategorys
	resp.Total = count
	resp.Message = "查询成功"

	return resp, nil

}
