package Product

type ProductService struct {
	Repo *ProductRepository
}

func NewProductService(repo *ProductRepository) *ProductService {
	return &ProductService{Repo: repo}
}

func (s *ProductService) CreateProduct(p *Product) error {
	// Burada istersen validasyon koyabilirsin
	return s.Repo.Create(p)
}

func (s *ProductService) UpdateProduct(p *Product) error {
	// Güncellemeden önce örneğin stok kontrolü yapılabilir
	return s.Repo.Update(p)
}

func (s *ProductService) DeleteProduct(id int) error {
	// Silmeden önce ilişkili sipariş var mı kontrol edilebilir
	return s.Repo.Delete(id)
}
