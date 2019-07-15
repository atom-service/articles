package provider

import (
	"context"

	"github.com/joho/godotenv"
	"github.com/yinxulai/grpc-services/articles/models"
	"github.com/yinxulai/grpc-services/articles/standard"
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

// Create Create
func (srv *Service) Create(ctx context.Context, req *standard.CreateRequest) (resp *standard.CreateResponse, err error) {
	var count uint64
	resp = new(standard.CreateResponse)
	// TODO: 检查分类是否存在
	if req.Article.OwnerCategory == 0 {
		resp.State = standard.State_PARAMS_INVALID
		resp.Message = "无效的分类"
		return resp, nil
	}

	err = countCategoryByIDStmt.GetContext(ctx, &count, map[string]interface{}{"ID": req.Article.OwnerCategory})
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
	_, err = insertArticleStmt.ExecContext(ctx, req.Article)
	if err != nil {
		resp.State = standard.State_DB_OPERATION_FATLURE
		resp.Message = err.Error()
		return resp, nil
	}

	resp.State = standard.State_SUCCESS
	resp.Message = "创建成功"

	return resp, nil
}

// QueryByID QueryByID
func (srv *Service) QueryByID(ctx context.Context, req *standard.QueryByIDRequest) (resp *standard.QueryByIDResponse, err error) {
	articles := []*models.Article{}
	resp = new(standard.QueryByIDResponse)

	if req.ID == 0 {
		resp.State = standard.State_PARAMS_INVALID
		resp.Message = "无效的 ID"
		return resp, nil
	}

	rows, err := queryArticleByIDStmt.QueryxContext(ctx, req)
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

// DeleteByID DeleteByID
func (srv *Service) DeleteByID(ctx context.Context, req *standard.DeleteByIDRequest) (resp *standard.DeleteByIDResponse, err error) {
	var count uint64
	resp = new(standard.DeleteByIDResponse)

	if req.ID == 0 {
		resp.State = standard.State_PARAMS_INVALID
		resp.Message = "无效的 ID"
		return resp, nil
	}

	err = countArticleByIDStmt.GetContext(ctx, &count, req)
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

	_, err = deleteArticleByIDStmt.ExecContext(ctx, req)
	if err != nil {
		resp.State = standard.State_DB_OPERATION_FATLURE
		resp.Message = err.Error()
		return resp, nil
	}

	resp.State = standard.State_SUCCESS
	resp.Message = "删除成功"

	return resp, nil
}

// UpdateByID UpdateByID
func (srv *Service) UpdateByID(ctx context.Context, req *standard.UpdateByIDRequest) (resp *standard.UpdateByIDResponse, err error) {
	// 检查是否存在该记录
	var count uint64
	resp = new(standard.UpdateByIDResponse)

	if req.ID == 0 {
		resp.State = standard.State_PARAMS_INVALID
		resp.Message = "无效的 ID"
		return resp, nil
	}

	err = countArticleByIDStmt.GetContext(ctx, &count, req)
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
	_, err = updateArticleByIDStmt.ExecContext(ctx, req.Data)
	if err != nil {
		resp.State = standard.State_DB_OPERATION_FATLURE
		resp.Message = err.Error()
		return resp, nil
	}

	resp.State = standard.State_SUCCESS
	resp.Message = "更新成功"

	return resp, nil
}

// QueryByOwner QueryByOwner
func (srv *Service) QueryByOwner(ctx context.Context, req *standard.QueryByOwnerRequest) (resp *standard.QueryByOwnerResponse, err error) {
	var count uint64
	articles := []*models.Article{}
	stdarticles := []*standard.Article{}
	resp = new(standard.QueryByOwnerResponse)
	// 查询 Owner 文章是否存在

	if req.Owner == 0 {
		resp.State = standard.State_PARAMS_INVALID
		resp.Message = "无效的 Owner"
		return resp, nil
	}

	// 查询记录总数
	err = countArticleByOwnerStmt.GetContext(ctx, &count, req)
	if err != nil {
		resp.State = standard.State_DB_OPERATION_FATLURE
		resp.Message = err.Error()
		return resp, nil
	}

	rows, err := queryArticleByOwnerStmt.QueryxContext(ctx, req)
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

// QueryByOwnerCategory QueryByOwnerCategory
func (srv *Service) QueryByOwnerCategory(ctx context.Context, req *standard.QueryByOwnerCategoryRequest) (resp *standard.QueryByOwnerCategoryResponse, err error) {
	// 检查是否存在该记录
	var count uint64
	articles := []*models.Article{}
	stdarticles := []*standard.Article{}
	resp = new(standard.QueryByOwnerCategoryResponse)

	// 简单检查一下分类
	if req.OwnerCategory == 0 {
		resp.State = standard.State_PARAMS_INVALID
		resp.Message = "无效的分类"
		return resp, nil
	}

	// 查询 Category 是否存在
	err = countCategoryByIDStmt.GetContext(ctx, &count, map[string]interface{}{"ID": req.OwnerCategory})
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
	err = countArticleByOwnerCategoryStmt.GetContext(ctx, &count, req)
	if err != nil {
		resp.State = standard.State_DB_OPERATION_FATLURE
		resp.Message = err.Error()
		return resp, nil
	}

	// 执行条件查询
	rows, err := queryArticleByOwnerCategoryStmt.QueryxContext(ctx, req)
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

	rows, err := queryLabelByIDStmt.QueryxContext(ctx, req)
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

	err = countLabelByIDStmt.GetContext(ctx, &count, req)
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
	_, err = updateLabelByIDStmt.ExecContext(ctx, req.Data)
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

	err = countLabelByIDStmt.GetContext(ctx, &count, req)
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

	_, err = deleteLabelByIDStmt.ExecContext(ctx, req)
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

	err = countArticleByIDStmt.GetContext(ctx, &count, map[string]interface{}{"ID": req.Owner})
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
	_, err = insertLabelStmt.ExecContext(ctx, req.Label)
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

	err = countArticleByIDStmt.GetContext(ctx, &count, map[string]interface{}{"ID": req.Owner})
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
	err = countLabelByOwnerStmt.GetContext(ctx, &count, req)
	if err != nil {
		resp.State = standard.State_DB_OPERATION_FATLURE
		resp.Message = err.Error()
		return resp, nil
	}

	rows, err := queryLabelByOwnerStmt.QueryxContext(ctx, req)
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
	_, err = insertCategoryStmt.ExecContext(ctx, req.Category)
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

	err = countCategoryByIDStmt.GetContext(ctx, &count, req)
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
	_, err = updateCategoryByIDStmt.ExecContext(ctx, req.Data)
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

	err = countCategoryByIDStmt.GetContext(ctx, &count, req)
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
	_, err = deleteCategoryByIDStmt.ExecContext(ctx, req)
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

	err = countCategoryByOwnerStmt.GetContext(ctx, &count, req)
	if err != nil {
		resp.State = standard.State_DB_OPERATION_FATLURE
		resp.Message = err.Error()
		return resp, nil
	}

	rows, err := queryCategoryByOwnerStmt.QueryxContext(ctx, req)
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
	err = countCategoryByOwnerCategoryStmt.GetContext(ctx, &count, req)
	if err != nil {
		resp.State = standard.State_DB_OPERATION_FATLURE
		resp.Message = err.Error()
		return resp, nil
	}

	rows, err := queryCategoryByOwnerCategoryStmt.QueryxContext(ctx, req)
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
